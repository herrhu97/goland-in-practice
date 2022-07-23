package runner

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"time"
)

/**
调度后台处理任务的程序
*/

type Runner struct {
	interrupt chan os.Signal

	complete chan error

	timeout <-chan time.Time

	tasks []func(int)

	concurrency chan struct{}
}

var ErrTimeout = errors.New("received timeout")

var ErrInterrupt = errors.New("received interrupt")

func New(d time.Duration, con int) *Runner {
	return &Runner{
		interrupt:   make(chan os.Signal, 1),
		complete:    make(chan error),
		timeout:     time.After(d),
		concurrency: make(chan struct{}, con),
	}
}

func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

func (r *Runner) run() error {
	var wg sync.WaitGroup
	wg.Add(len(r.tasks))
	for id, task := range r.tasks {
		if r.gotInterrupt() {
			return ErrInterrupt
		}

		go func(id int) {
			select {
			case r.concurrency <- struct{}{}:
				task(id)
				wg.Done()
				<-r.concurrency
			}
		}(id)
	}
	wg.Wait()
	return nil
}

func (r *Runner) gotInterrupt() bool {
	select {
	case <-r.interrupt:
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}

func (r *Runner) Start() error {
	signal.Notify(r.interrupt, os.Interrupt)

	go func() {
		r.complete <- r.run()
	}()

	select {
	case err := <-r.complete:
		return err

	case <-r.timeout:
		return ErrTimeout
	}
}
