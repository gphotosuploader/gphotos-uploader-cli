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

type DeletionJob struct {
	uploadedFileURL string
	localFilePath   string
	typedMedia      filetypes.TypedMedia
}

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

func CloseDeletionsChan() { close(deletionsChan) }

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
