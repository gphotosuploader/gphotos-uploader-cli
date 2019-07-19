package upload

import (
	"log"
	"os"

	"github.com/juju/errors"
	"github.com/nmrshll/gphotos-uploader-cli/filetypes"
)

var (
	deletionsChan = make(chan DeletionJob)
)

// DeletionJob represents an object to be deleted from local repository
type DeletionJob struct {
	uploadedFileURL string
	localFilePath   string
	typedMedia      filetypes.TypedMedia
}

// TODO: create a new type DeletionQueue. The rest are methods of it

// QueueDeletionJob adds an object to be deleted to an asynchronous queue.
// It checks that all the parameters are set or return error otherwise.
func QueueDeletionJob(deletionJob DeletionJob) error {
	// check params
	{
		if deletionJob.uploadedFileURL == "" {
			return errors.NewNotValid(nil, "missing uploadedFileURL for deletionJob")
		}
		if deletionJob.localFilePath == "" {
			return errors.NewNotValid(nil, "missing localFilePath for deletionJob")
		}
		if deletionJob.typedMedia == nil {
			return errors.NewNotValid(nil, "missing typedMedia for deletionJob")
		}
	}

	deletionsChan <- deletionJob
	return nil
}

// CloseDeletionChan close the channel used for removing objects
func CloseDeletionsChan() { close(deletionsChan) }


// StartDeletionsWorker set up channels and start concurrent deletions
// deletionsChan will receive DeletionJob structs and delete the object
// from local repository, it will signal doneDeleting when removal is done
func StartDeletionsWorker() (doneDeleting chan struct{}) {
	doneDeleting = make(chan struct{})
	go func() {
		for deletionJob := range deletionsChan {
			_ = deletionJob.deleteIfCorrectlyUploaded()
		}
		doneDeleting <- struct{}{}
	}()
	return doneDeleting
}

func (deletionJob *DeletionJob) deleteIfCorrectlyUploaded() error {
	isCorrectlyUploaded, err := deletionJob.typedMedia.IsCorrectlyUploaded(deletionJob.uploadedFileURL, deletionJob.localFilePath)
	if err != nil {
		log.Printf("%s. Won't delete\n", err)
		return err
	}

	if isCorrectlyUploaded {
		log.Printf("uploaded file %s was checked for integrity. Will now delete.\n", deletionJob.localFilePath)
		if err = os.Remove(deletionJob.localFilePath); err != nil {
			log.Println("failed deleting file")
		}

		//if err = RemoveAsAlreadyUploaded(deletionJob.localFilePath); err != nil {
		//	log.Printf("Failed to remove from DB: %s", err)
		//}
		return err
	}

	log.Println("not the same image. Won't delete")
	return err
}
