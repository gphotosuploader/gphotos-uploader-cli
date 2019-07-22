package upload

import (
	"fmt"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"
	"github.com/juju/errors"
	"github.com/nmrshll/go-cp"
	"github.com/nmrshll/gphotos-uploader-cli/config"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/completeduploads"
	"github.com/nmrshll/gphotos-uploader-cli/filetypes"
	"github.com/nmrshll/gphotos-uploader-cli/utils/filesystem"
	"log"
	"os"
	"path/filepath"
)

// Job represents a job to upload all photos in a folder
type Job struct {
	client          *gphotos.Client
	trackingService *completeduploads.Service

	SourceFolder      string
	MakeAlbums        config.MakeAlbums
	DeleteAfterUpload bool
	UploadVideos      bool
	IncludePatterns   []string
	ExcludePatterns   []string
}

// NewFolderUploadJob creates a Job based on the submitted data
func NewFolderUploadJob(client *gphotos.Client, trackingService *completeduploads.Service, cfg *config.FolderUploadJob) *Job {
	return &Job{
		trackingService: trackingService,
		client:          client,

		SourceFolder:      cfg.SourceFolder,
		MakeAlbums:        cfg.MakeAlbums,
		DeleteAfterUpload: cfg.DeleteAfterUpload,
		UploadVideos:      cfg.UploadVideos,
		IncludePatterns:   cfg.IncludePatterns,
		ExcludePatterns:   cfg.ExcludePatterns,
	}
}

// ScanFolder uploads folder
func (job *Job) ScanFolder(uploadChan chan<- *Item) error {
	folderAbsolutePath, err := cp.AbsolutePath(job.SourceFolder)
	if err != nil {
		return err
	}

	if !filesystem.IsDir(folderAbsolutePath) {
		return fmt.Errorf("%s is not a folder", folderAbsolutePath)
	}

	filter := NewFilter(job.IncludePatterns, job.ExcludePatterns, job.UploadVideos)

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

		typedMedia, err := filetypes.NewTypedMedia(fp)
		if err != nil {
			log.Println(errors.Annotatef(err, "failed creating new TypedMedia from fp"))
			return nil
		}

		// calculate Album Name from the folder name
		album := filepath.Base(filepath.Dir(fp))

		// set file upload options depending on folder upload options
		var uploadItem = &Item{
			client:          job.client,
			path:            fp,
			typedMedia:      typedMedia,
			album:           album,
			deleteOnSuccess: job.DeleteAfterUpload,
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
