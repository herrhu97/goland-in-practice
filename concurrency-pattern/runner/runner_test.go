package runner

import (
	"log"
	"os"
	"testing"
	"time"
)

const timeout = 6 * time.Second

func TestRunner(t *testing.T) {
	log.Println("Starting work.")

	r := New(timeout, 3)

	r.Add(createTask(), createTask(), createTask(), createTask(), createTask(), createTask())

	if err := r.Start(); err != nil {
		switch err {
		case ErrTimeout:
			log.Println("Terminating due to timeout.")
			os.Exit(1)
		case ErrInterrupt:
			log.Println("Terminating due to interrupt.")
			os.Exit(2)
		}
	}

	log.Println("Process ended.")
}

func createTask() func(int) {
	return func(id int) {
		log.Printf("Processor - Task #%d.", id)
		log.Printf("Task #%d, time before:%s", id, time.Now().Format("2006-01-02 15:04:05.000"))
		time.Sleep(time.Duration(id) * time.Second)
		log.Printf("Task #%d, time after:%s", id, time.Now().Format("2006-01-02 15:04:05.000"))
	}
}
