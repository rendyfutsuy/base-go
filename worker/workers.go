package worker

import (
	"log"
	// 💡 2. Import the packages containing the usecase INTERFACES, not the implementation folders.
	// ..
)

// UsecaseRegistry holds all the usecase interfaces that the worker might need.
type UsecaseRegistry struct {
	// Add other usecase interfaces here as needed
}

// Worker executes jobs from the queue.
type Worker struct {
	ID         int
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
	usecases   UsecaseRegistry // Worker has access to all usecases
}

func NewWorker(workerPool chan chan Job, id int, usecases UsecaseRegistry) Worker {
	return Worker{
		ID:         id,
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
		usecases:   usecases,
	}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				log.Printf("Worker %d: received job %s of type %s\n", w.ID, job.ID, job.Type)
				w.processJob(job)
			case <-w.quit:
				return
			}
		}
	}()
}

// processJob is the worker's router. It delegates the job to the correct usecase.
func (w Worker) processJob(job Job) {
	// TBA
}
