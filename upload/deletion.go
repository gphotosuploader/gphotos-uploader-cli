package upload

import (
	"log"
	"os"

	"github.com/juju/errors"
)

var (
	deletionQueue = make(chan DeletionJob, 25)
)

// DeletionJob represents an object to be deleted from local repository.
type DeletionJob struct {
	ObjectURL  string
	ObjectPath string
}

// Enqueue adds an object to be deleted to an asynchronous queue.
func (j *DeletionJob) Enqueue() error {
	deletionQueue <- *j
	return nil
}

func (j *DeletionJob) deleteIfCorrectlyUploaded() error {
	if j.ObjectURL == "" {
		return errors.NewNotValid(nil, "file was not correctly uploaded, skipping deletion")
	}

	err := os.Remove(j.ObjectPath)
	return err
}

// DeletionQueue is an async queue to process deletion requests.
// Requests will receive DeletionJob structs to be deleted from local repository.
// DoneChannel will be signalled when removal is done.
type DeletionQueue struct {
	Requests    chan DeletionJob
	DoneChannel chan bool
}

// NewDeletionQueue returns a new DeletionQueue object.
func NewDeletionQueue() *DeletionQueue {
	return &DeletionQueue{
		Requests:    deletionQueue,
		DoneChannel: make(chan bool),
	}
}

// StartWorkers start concurrent deletion workers.
func (q *DeletionQueue) StartWorkers() {
	go func() {
		for job := range q.Requests {
			log.Printf("Deletion worker: Processing deletion request for: file=%s", job.ObjectPath)
			err := job.deleteIfCorrectlyUploaded()
			if err != nil {
				log.Printf("File has not been deleted: file=%s, err=%v", job.ObjectPath, err)
			}
		}
		q.DoneChannel <- true
	}()
}

// Close closes the queue to avoid new jobs are added.
func (q *DeletionQueue) Close() {
	close(q.Requests)
}

// WaitForWorkers waits for deletions to be completed.
// It waits until quit channel is read.
func (q *DeletionQueue) WaitForWorkers() {
	<-q.DoneChannel
}
