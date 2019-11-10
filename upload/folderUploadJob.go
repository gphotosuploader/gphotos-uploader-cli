package upload

import (
	"os"
	"path/filepath"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"

	"github.com/gphotosuploader/gphotos-uploader-cli/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/log"
	"github.com/gphotosuploader/gphotos-uploader-cli/utils/filesystem"
)

// Job represents a job to upload all photos from the specified folder
type Job struct {
	gPhotos     *gphotos.Client
	fileTracker app.FileTracker

	sourceFolder string
	options      *JobOptions
}

// JobOptions represents all the options that a job can have
type JobOptions struct {
	createAlbum        bool
	createAlbumBasedOn string
	deleteAfterUpload  bool
	filter             *Filter
}

// NewJobOptions create a jobOptions based on the submitted / validated data
func NewJobOptions(createAlbum bool, createAlbumBasedOn string, deleteAfterUpload bool, filter *Filter) *JobOptions {
	return &JobOptions{
		createAlbum:        createAlbum,
		createAlbumBasedOn: createAlbumBasedOn,
		deleteAfterUpload:  deleteAfterUpload,
		filter:             filter,
	}
}

// NewFolderUploadJob creates a job based on the submitted data
func NewFolderUploadJob(gPhotos *gphotos.Client, fileTracker app.FileTracker, fp string, opt *JobOptions) *Job {
	return &Job{
		fileTracker: fileTracker,
		gPhotos:     gPhotos,

		sourceFolder: fp,
		options:      opt,
	}
}

// ScanFolder return the list of Items{} to be uploaded. It scans the folder and skip
// non allowed files (includePatterns & excludePattens).
func (job *Job) ScanFolder(logger log.Logger) ([]Item, error) {
	var result []Item
	err := filepath.Walk(job.sourceFolder, job.createUploadItemListFn(&result, job.options.filter, logger))
	return result, err
}

func (job *Job) createUploadItemListFn(items *[]Item, filter *Filter, logger log.Logger) filepath.WalkFunc {
	return func(fp string, fi os.FileInfo, errP error) error {
		if fi == nil || fi.IsDir() {
			return nil
		}

		// check if the item should be uploaded given the include and exclude patterns in the
		// configuration file. It uses relative path from the source folder path to facilitate
		// then set up of includePatterns and excludePatterns.
		relativePath := filesystem.RelativePath(job.sourceFolder, fp)
		if !filter.IsAllowed(relativePath) {
			logger.Infof("Not allowed by config: %s: skipping file...", fp)
			return nil
		}

		// check completed uploads db for previous uploads
		isAlreadyUploaded, err := job.fileTracker.IsAlreadyUploaded(fp)
		if err != nil {
			logger.Error(err)
		} else if isAlreadyUploaded {
			logger.Debugf("Already uploaded: %s: skipping file...", fp)
			return nil
		}

		logger.Infof("Adding item to upload: %s", fp)

		// set file upload options depending on folder upload options
		*items = append(*items, Item{
			gPhotos:         job.gPhotos,
			path:            fp,
			album:           job.albumID(relativePath, logger),
			deleteOnSuccess: job.options.deleteAfterUpload,
		})
		return nil
	}
}
