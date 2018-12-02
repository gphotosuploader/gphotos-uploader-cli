package fileshandling

import (
	"fmt"
	imageLib "image"
	"log"
	"strconv"

	// register decoders for jpeg and png
	_ "image/jpeg"
	_ "image/png"

	"net/http"
	"os"
	"strings"

	"github.com/Nr90/imgsim"
	"github.com/nmrshll/gphotos-uploader-cli/util"
	"github.com/steakknife/hamming"
	"github.com/syndtr/goleveldb/leveldb"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

var (
	imageExtensions = []string{".jpg", ".jpeg", ".png", ".webp"}
	deletionsChan   = make(chan DeletionJob)
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
			err := CheckUploadedAndDeleteLocal(deletionJob.uploadedMediaItem, deletionJob.localFilePath, db)
			if err != nil {
				fmt.Printf("%s. Won't delete", err)
			}
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
	log.Printf("Removing file from upload DB: %s", filePath, val)
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

func getImageHash(img imageLib.Image) string {
	return imgsim.DifferenceHash(img).String()
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

// CheckUploadedAndDeleteLocal checks that the image that was uploaded is visually similar to the local one, before deleting the local one
func CheckUploadedAndDeleteLocal(uploadedMediaItem *photoslibrary.MediaItem, localImgPath string, db *leveldb.DB) error {
	if !HasImageExtension(localImgPath) {
		return fmt.Errorf("%s doesn't have an image extension", localImgPath)
	}

	// compare uploaded image and local one
	upImg, err := imageFromURL(uploadedMediaItem.BaseUrl)
	if err != nil {
		return fmt.Errorf("failed getting image from URL")
	}
	localImg, err := imageFromPath(localImgPath)
	if err != nil {
		return fmt.Errorf("failed loading local image from path")
	}

	if !isSameImage(upImg, localImg) {
		log.Println("not the same image. Won't delete")
	} else {
		log.Printf("uploaded file %s was checked for integrity. Will now delete.", localImgPath)
		if err = os.Remove(localImgPath); err != nil {
			log.Println("delete failed")
		} else {
			if err = removeFromDB(localImgPath, db); err != nil {
				log.Printf("Failed to remove from DB: %s", err)
			}
		}
	}
	return nil
}
