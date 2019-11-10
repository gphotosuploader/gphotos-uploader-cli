package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	uploadVideos       bool
	includePatterns    []string
	excludePatterns    []string
}

// NewJobOptions create a jobOptions based on the submitted / validated data
func NewJobOptions(createAlbum bool, createAlbumBasedOn string, deleteAfterUpload bool, uploadVideos bool, includePatterns []string, excludePatterns []string) *JobOptions {
	return &JobOptions{
		createAlbum:        createAlbum,
		createAlbumBasedOn: createAlbumBasedOn,
		deleteAfterUpload:  deleteAfterUpload,
		uploadVideos:       uploadVideos,
		includePatterns:    includePatterns,
		excludePatterns:    excludePatterns,
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
	filter := NewFilter(job.options.includePatterns, job.options.excludePatterns, job.options.uploadVideos)

	var result []Item
	err := filepath.Walk(job.sourceFolder, job.createUploadItemListFn(&result, filter, logger))
	return result, err
}

func (job *Job) createUploadItemListFn(items *[]Item, filter *Filter, logger log.Logger) filepath.WalkFunc {
	return func(fp string, fi os.FileInfo, errP error) error {
		if fi == nil {
			logger.Fatalf("error scanning: folder=%s, err=FileInfo is nil", fp)
			return nil
		}

		// avoid processing folders
		if fi.IsDir() {
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

		logger.Infof("Adding new item to upload list: item=%s, rel=%s", fp, relativePath)

		// calculate Album from the folder name, we create if it's not exists
		var albumID string
		if job.options.createAlbum {
			albumID, err = job.createAlbum(relativePath)
			if err != nil {
				logger.Error(err)
			}
		}

		// set file upload options depending on folder upload options
		var uploadItem = Item{
			gPhotos:         job.gPhotos,
			path:            fp,
			album:           albumID,
			deleteOnSuccess: job.options.deleteAfterUpload,
		}

		*items = append(*items, uploadItem)
		return nil
	}
}

// createAlbum returns the ID of an album with the specified name or error if fails.
// If the album didn't exist, it's created thanks to GetOrCreateAlbumByName().
func (job *Job) createAlbum(path string) (string, error) {
	var name string
	switch job.options.createAlbumBasedOn {
	case "folderPath":
		name = strings.ReplaceAll(filepath.Dir(path), "/", "_")
	case "folderName":
	default:
		name = filepath.Base(filepath.Dir(path))
	}

	if name == "" {
		return "", nil
	}

	// get album ID from Google Photos API
	album, err := job.gPhotos.GetOrCreateAlbumByName(name)
	if err != nil {
		return "", fmt.Errorf("album creation failed: name=%s, error=%s", name, err)
	}
	return album.Id, nil
}
