package worker

import "log"

// Dispatcher manages the pool of workers.
type Dispatcher struct {
	WorkerPool chan chan Job
	maxWorkers int
	usecases   UsecaseRegistry
}

// NewDispatcher creates a new dispatcher. You will call this in your main.go.
func NewDispatcher(maxWorkers int, usecases UsecaseRegistry) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		WorkerPool: pool,
		maxWorkers: maxWorkers,
		usecases:   usecases,
	}
}

// Run starts the dispatcher and all the workers.
func (d *Dispatcher) Run() {
	JobQueue = make(chan Job, 100) // Initialize the global queue
	for i := 1; i <= d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, i, d.usecases)
		worker.Start()
	}
	go d.dispatch()
	log.Println("Background worker system is running.")
}

// dispatch listens on the global JobQueue and sends jobs to available workers.
func (d *Dispatcher) dispatch() {
	for job := range JobQueue {
		go func(job Job) {
			jobChannel := <-d.WorkerPool
			jobChannel <- job
		}(job)
	}
}
