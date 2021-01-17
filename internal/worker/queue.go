package worker

import (
	"fmt"
	"sync"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
)

// Job - interface for job processing
type Job interface {
	Process() error
	ID() string
}

// JobResult - the jobResults of a processed Job
type JobResult struct {
	ID      string
	Message string
	Err     error
}

// Worker - the worker threads that actually process the jobs
type Worker struct {
	id               string
	done             *sync.WaitGroup
	readyPool        chan chan Job
	assignedJobQueue chan Job

	jobResults chan JobResult
	quit       chan bool

	logger log.Logger
}

// JobQueue - a queue for enqueueing jobs to be processed
type JobQueue struct {
	internalQueue     chan Job
	readyPool         chan chan Job
	workers           []*Worker
	dispatcherStopped sync.WaitGroup
	workersStopped    *sync.WaitGroup
	jobResults        chan JobResult
	quit              chan bool
}

// NewJobQueue - creates a new job queue
func NewJobQueue(maxWorkers int, logger log.Logger) *JobQueue {
	workersStopped := sync.WaitGroup{}
	readyPool := make(chan chan Job, maxWorkers)

	// we need to ensure that the results channel is big enough to fulfill workers needs
	jobResults := make(chan JobResult, maxWorkers*10)

	// create the pool of workers
	workers := make([]*Worker, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		workers[i] = NewWorker(fmt.Sprintf("#%d", i+1), readyPool, jobResults, &workersStopped, logger)
	}
	return &JobQueue{
		internalQueue:     make(chan Job),
		readyPool:         readyPool,
		workers:           workers,
		dispatcherStopped: sync.WaitGroup{},
		workersStopped:    &workersStopped,
		jobResults:        jobResults,
		quit:              make(chan bool),
	}
}

func (q *JobQueue) ChanJobResults() chan JobResult {
	return q.jobResults
}

// Start - starts the worker routines and dispatcher routine
func (q *JobQueue) Start() {
	for i := 0; i < len(q.workers); i++ {
		q.workers[i].Start()
	}
	go q.dispatch()
}

// Stop - stops the workers and dispatcher routine
func (q *JobQueue) Stop() {
	q.quit <- true
	q.dispatcherStopped.Wait()
}

func (q *JobQueue) dispatch() {
	q.dispatcherStopped.Add(1)
	for {
		select {
		case job := <-q.internalQueue: // We got something in on our queue
			workerChannel := <-q.readyPool // Check out an available worker
			workerChannel <- job           // Send the request to the channel
		case <-q.quit:
			for i := 0; i < len(q.workers); i++ {
				q.workers[i].Stop()
			}
			q.workersStopped.Wait()
			q.dispatcherStopped.Done()
			return
		}
	}
}

// Submit - adds a new job to be processed, uses a subroutine to avoid deadlock when the queue is full
func (q *JobQueue) Submit(job Job) {
	go func(job Job) { q.internalQueue <- job }(job)
}

// NewWorker - creates a new worker
func NewWorker(id string, readyPool chan chan Job, result chan JobResult, done *sync.WaitGroup, logger log.Logger) *Worker {
	return &Worker{
		id:               id,
		done:             done,
		readyPool:        readyPool,
		assignedJobQueue: make(chan Job),
		jobResults:       result,
		quit:             make(chan bool),
		logger:           logger,
	}
}

// Start - begins the job processing loop for the worker
func (w *Worker) Start() {
	go func() {
		w.logger.Debugf("Worker %s is starting", w.id)
		w.done.Add(1)
		for {
			w.readyPool <- w.assignedJobQueue // check the job queue in
			select {
			case job := <-w.assignedJobQueue: // see if anything has been assigned to the queue
				w.logger.Debugf("Worker %s processing: %s", w.id, job.ID())

				r := JobResult{
					ID:      job.ID(),
					Message: "processed successfully",
					Err:     job.Process(),
				}

				if r.Err != nil {
					r.Message = "processed with errors"
					w.logger.Error(r.Err)
				}

				// send the jobResults of the processed Job
				w.jobResults <- r
			case <-w.quit:
				w.done.Done()
				return
			}
		}
	}()
}

// Stop - stops the worker
func (w *Worker) Stop() {
	w.logger.Debugf("Worker %s is stopping", w.id)
	w.quit <- true
}
