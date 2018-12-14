package fileshandling

import (
	"fmt"
	imageLib "image"

	// register decoders for jpeg and png
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"net/http"
	"os"

	"github.com/Nr90/imgsim"
	"github.com/palantir/stacktrace"
	"github.com/steakknife/hamming"
	"github.com/syndtr/goleveldb/leveldb"
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

func StartDeletionsWorker(db *leveldb.DB) (doneDeleting chan struct{}) {
	doneDeleting = make(chan struct{})
	go func() {
		for deletionJob := range deletionsChan {
			deletionJob.deleteIfCorrectlyUploaded()
		}
		doneDeleting <- struct{}{}
	}()
	return doneDeleting
}

func IsUploadedPrev(filePath string, db *leveldb.DB) (bool, error) {
	isUploaded := false

	// look for previous upload in cache
	val, err := db.Get([]byte(filePath), nil)
	if err == nil {
		// value found, try to split mtime and hash
		parts := strings.Split(string(val[:]), "|")
		cacheMtime := int64(0)
		cacheHash := ""
		if len(parts) > 1 {
			cacheMtime, err = strconv.ParseInt(parts[0], 10, 64)
			cacheHash = parts[1]
		} else {
			cacheHash = parts[0]
		}
		// check mtime first
		if err == nil && cacheMtime != 0 {
			fileMtime, err := util.GetMTime(filePath)
			if err == nil && fileMtime.Unix() == cacheMtime {
				isUploaded = true
				//log.Printf("%s mtime matched %i", filePath, cacheMtime)
			}
		}
		// mtime is different, check hash
		if !isUploaded {
			localImg, err := imageFromPath(filePath)
			if err != nil {
				err = fmt.Errorf("failed loading local image from path")
			} else {
				localHash := getImageHash(localImg)
				isUploaded = isSameHash(cacheHash, localHash)
				if isUploaded {
					//log.Printf("%s hash match %s", filePath, cacheHash)
					// update db mtime
					err = MarkUploaded(filePath, db)
				}
			}
		}
	} else if strings.Contains(err.Error(), "not found") {
		err = nil
	}

	return isUploaded, err
}

func MarkUploaded(filePath string, db *leveldb.DB) error {
	localImg, err := imageFromPath(filePath)
	if err != nil {
		return fmt.Errorf("failed loading local image from path")
	}
	mtime, err := util.GetMTime(filePath)
	if err != nil {
		return fmt.Errorf("failed getting local image mtime")
	}
	val := strconv.FormatInt(mtime.Unix(), 10) + "|" + getImageHash(localImg)

	log.Printf("Marking file as uploaded: %s", filePath, val)
	err = db.Put([]byte(filePath), []byte(val), nil)

	return err
}

func removeFromDB(filePath string, db *leveldb.DB) error {
	log.Printf("Removing file from upload DB: %s", filePath)
	err := db.Delete([]byte(filePath), nil)

	return err
}

func HasImageExtension(path string) bool {
	for _, ext := range imageExtensions {
		if strings.HasSuffix(strings.ToLower(path), ext) {
			return true
		}
	}
	return false
}

func imageFromPath(filePath string) (imageLib.Image, error) {
	reader, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	img, _, err := imageLib.Decode(reader)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func imageFromURL(URL string) (imageLib.Image, error) {
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("expected http status 200, got %d", res.StatusCode)
	}

	img, _, err := imageLib.Decode(res.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func isSameImage(upImg, localImg imageLib.Image) bool {
	upDHash := getImageHash(upImg)
	localDHash := getImageHash(localImg)

	return isSameHash(upDHash, localDHash)
}

func isSameHash(upDHash, localDHash string) bool {
	if len(upDHash) != len(localDHash) {
		return false
	}
	hammingDistance := hamming.Strings(upDHash, localDHash)

	if hammingDistance < len(upDHash)/16 {
		return true
	}
	return false
}

// isImageCorrectlyUploaded checks that the image that was uploaded is visually similar to the local one, before deleting the local one
func isImageCorrectlyUploaded(uploadedMediaItem *photoslibrary.MediaItem, localImgPath string) (bool, error) {
	// TODO: add sameness check for videos (use file hash) and delete if same
	if !IsImage(localImgPath) {
		return false, fmt.Errorf("%s is not an image. won't delete local file", localImgPath)
	}

	// compare uploaded image and local one
	upImg, err := imageFromURL(uploadedMediaItem.BaseUrl)
	if err != nil {
		return false, stacktrace.Propagate(err, "failed getting image from URL")
	}
	localImg, err := imageFromPath(localImgPath)
	if err != nil {
		return false, stacktrace.Propagate(err, "failed loading local image from path")
	}

	if isSameImage(upImg, localImg) {
		return true, nil
	}

	return false, nil
}

func (deletionJob *DeletionJob) deleteIfCorrectlyUploaded() {
	isImageCorrectlyUploaded, err := isImageCorrectlyUploaded(deletionJob.uploadedMediaItem, deletionJob.localFilePath)
	if err != nil {
		fmt.Printf("%s. Won't delete\n", err)
		return
	}

	if isImageCorrectlyUploaded {
		fmt.Printf("uploaded file %s was checked for integrity. Will now delete.\n", deletionJob.localFilePath)
		if err = os.Remove(deletionJob.localFilePath); err != nil {
			fmt.Println("failed deleting file")
		}
		return
	} else {
		fmt.Println("not the same image. Won't delete")
		return
	}
}
