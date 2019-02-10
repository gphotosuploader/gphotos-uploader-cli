package upload

import (
	"log"
	"os"

	"github.com/nmrshll/gphotos-uploader-cli/filetypes"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

var (
	deletionsChan = make(chan DeletionJob)
)

type DeletionJob struct {
	uploadedMediaItem *photoslibrary.MediaItem
	localFilePath     string
}

func QueueDeletionJob(uploadedMediaItem *photoslibrary.MediaItem, localImgPath string) {
	deletionsChan <- DeletionJob{uploadedMediaItem, localImgPath}
}

func CloseDeletionsChan() { close(deletionsChan) }

func StartDeletionsWorker() (doneDeleting chan struct{}) {
	doneDeleting = make(chan struct{})
	go func() {
		for deletionJob := range deletionsChan {
			_ = deletionJob.deleteImageIfCorrectlyUploaded()
		}
		doneDeleting <- struct{}{}
	}()
	return doneDeleting
}

func (deletionJob *DeletionJob) deleteImageIfCorrectlyUploaded() error {
	isImageCorrectlyUploaded, err := filetypes.IsImageCorrectlyUploaded(deletionJob.uploadedMediaItem, deletionJob.localFilePath)
	if err != nil {
		log.Printf("%s. Won't delete\n", err)
		return err
	}

	if isImageCorrectlyUploaded {
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
