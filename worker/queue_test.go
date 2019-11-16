package worker

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/log"
)

type TestJob struct {
	id     int
	opsPtr *uint64
}

func (t *TestJob) ID() string { return fmt.Sprintf("Worker #%d", t.id) }

func (t *TestJob) Process() error {
	for i := 0; i < 1000; i++ {
		atomic.AddUint64(t.opsPtr, 1)
	}
	return nil
}

func TestQueue(t *testing.T) {
	var testData = []struct {
		numberOfWorkers int
		numberOfJobs    int
		want            uint64
	}{
		{numberOfWorkers: 1, numberOfJobs: 1, want: 1000},
		{numberOfWorkers: 1, numberOfJobs: 5, want: 5000},
		{numberOfWorkers: 5, numberOfJobs: 5, want: 5000},
		{numberOfWorkers: 5, numberOfJobs: 50, want: 50000},
	}

	var logger = &log.DiscardLogger{}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Workers[%d]_Jobs[%d]", tt.numberOfWorkers, tt.numberOfJobs), func(t *testing.T) {
			var ops uint64

			queue := NewJobQueue(tt.numberOfWorkers, logger)
			queue.Start()

			for i := 0; i < tt.numberOfJobs; i++ {
				queue.Submit(&TestJob{id: i, opsPtr: &ops})
			}

			queue.Stop()

			if ops != tt.want {
				t.Errorf("test failed: want=%d, got=%d", tt.want, ops)
			}
		})
	}

}
