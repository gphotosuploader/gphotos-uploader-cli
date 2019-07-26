package upload

import (
	"fmt"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"
	"github.com/nmrshll/go-cp"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/completeduploads"
	"github.com/nmrshll/gphotos-uploader-cli/utils/filesystem"
	"log"
	"os"
	"path/filepath"
)

// job represents a job to upload all photos from the specified folder
type job struct {
	client          *gphotos.Client
	trackingService *completeduploads.Service

	sourceFolder string
	*jobOptions
}

// jobOptions represents all the options that a job can have
type jobOptions struct {
	createAlbum       bool
	deleteAfterUpload bool
	uploadVideos      bool
	includePatterns   []string
	excludePatterns   []string
}

// NewJobOptions create a jobOptions based on the submitted / validated data
func NewJobOptions(createAlbum bool, deleteAfterUpload bool, uploadVideos bool, includePatterns []string, excludePatterns []string) *jobOptions {
	return &jobOptions{
		createAlbum:       createAlbum,
		deleteAfterUpload: deleteAfterUpload,
		uploadVideos:      uploadVideos,
		includePatterns:   includePatterns,
		excludePatterns:   excludePatterns,
	}
}

// NewFolderUploadJob creates a job based on the submitted data
func NewFolderUploadJob(client *gphotos.Client, trackingService *completeduploads.Service, fp string, opt *jobOptions) *job {
	return &job{
		trackingService: trackingService,
		client:          client,

		sourceFolder: fp,
		jobOptions:   opt,
	}
}

// ScanFolder uploads folder
func (job *job) ScanFolder(uploadChan chan<- *Item) error {
	folderAbsolutePath, err := cp.AbsolutePath(job.sourceFolder)
	if err != nil {
		return err
	}

	if !filesystem.IsDir(folderAbsolutePath) {
		return fmt.Errorf("%s is not a folder", folderAbsolutePath)
	}

	filter := NewFilter(job.includePatterns, job.excludePatterns, job.uploadVideos)

	// dirs are walked depth-first.   These vars hold the active album
	// default empty album for makeAlbums.enabled = false
	errW := filepath.Walk(folderAbsolutePath, func(fp string, fi os.FileInfo, errP error) error {
		// log.Printf("ScanFolder.Walk: %v, fi: %v, err: %v\n", fp, fi, err)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error for %v: %v\n", fp, err)
			return nil
		}
		if fi == nil {
			_, _ = fmt.Fprintf(os.Stderr, "error for %v: FileInfo is nil\n", fp)
			return nil
		}

		// check if the item should be uploaded (it's a file and it's not exclude
		if !filter.IsAllowed(fp) {
			if fi.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// only files are allowed
		if !filesystem.IsFile(fp) {
			return nil
		}

		// check upload db for previous uploads
		isAlreadyUploaded, err := job.trackingService.IsAlreadyUploaded(fp)
		if err != nil {
			log.Println(err)
		} else if isAlreadyUploaded {
			log.Printf("already uploaded: %s: skipping file...\n", fp)
			return nil
		}

		// calculate Album Name from the folder name
		var album string
		if job.createAlbum {
			album = filepath.Base(filepath.Dir(fp))
		}

		// TODO: Fix issue #25 - Removal of GIF & Videos is broken: https://github.com/nmrshll/gphotos-uploader-cli/issues/25
		// v0.4.0: Disable all files removal until we fix the issue properly
		if job.deleteAfterUpload {
			log.Printf("[WARNING] Removal of local files has been disabled. (See issue #25 https://github.com/nmrshll/gphotos-uploader-cli/issues/25")
		}

		// set file upload options depending on folder upload options
		var uploadItem = &Item{
			client:          job.client,
			path:            fp,
			album:           album,
			deleteOnSuccess: false, // TODO: Fix issue #25 - Removal of GIF & Videos is broken: https://github.com/nmrshll/gphotos-uploader-cli/issues/25
		}

		// finally, add the file upload to the queue
		uploadChan <- uploadItem

		return nil
	})
	if errW != nil {
		log.Printf("walk error [%v]", errW)
	}

	return nil
}
