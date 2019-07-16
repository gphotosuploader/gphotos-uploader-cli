package upload

import (
	"fmt"
	"github.com/nmrshll/go-cp"
	"github.com/nmrshll/gphotos-uploader-cli/filter"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/noserver-gphotos"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"github.com/juju/errors"
	"github.com/nmrshll/gphotos-uploader-cli/config"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/completeduploads"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/tokenstore"
	"github.com/nmrshll/gphotos-uploader-cli/filetypes"
	"github.com/nmrshll/gphotos-uploader-cli/utils/filesystem"
)

// FolderUploadJob represents a job to upload all photos in a folder
type FolderUploadJob struct {
	*config.FolderUploadJob
	uploaderConfigAPICredentials *config.APIAppCredentials
	gphotosClient                *gphotos.Client
	completedUploads             *completeduploads.Service
}

// NewFolderUploadJob creates a FolderUploadJob based on the submitted data
func NewFolderUploadJob(configFolderUploadJob *config.FolderUploadJob, completedUploads *completeduploads.Service, uploaderConfigAPICredentials *config.APIAppCredentials, tokenManagerService *tokenstore.Service) *FolderUploadJob {
	// check args
	{
		if completedUploads == nil {
			log.Fatalf("completedUploadsService can't be nil")
		}
		if uploaderConfigAPICredentials == nil {
			log.Fatalf("uploaderConfigAPICredentials can't be nil")
		}
	}

	folderUploadJob := &FolderUploadJob{
		FolderUploadJob:              configFolderUploadJob,
		completedUploads:             completedUploads,
		uploaderConfigAPICredentials: uploaderConfigAPICredentials,
	}

	gphotosClient, err := authenticate(tokenManagerService, folderUploadJob)
	if err != nil {
		log.Fatal(err)
	}
	folderUploadJob.gphotosClient = gphotosClient

	return folderUploadJob
}

func authenticate(tkm *tokenstore.Service, folderUploadJob *FolderUploadJob) (*gphotos.Client, error) {
	// try to load token from keyring

	token, err := tkm.RetrieveToken(folderUploadJob.Account)
	if err == nil && token != nil { // if error ignore and skip
		// if found create client from token
		gphotosClient, err := gphotos.NewClient(gphotos.FromToken(config.OAuthConfig(folderUploadJob.uploaderConfigAPICredentials), token))
		if err == nil && gphotosClient != nil { // if error ignore and skip
			return gphotosClient, nil
		}
	}

	// else authenticate again to grab a new token
	log.Println(color.CyanString(fmt.Sprintf("Need to log login into account %s", folderUploadJob.Account)))
	time.Sleep(1200 * time.Millisecond)
	gphotosClient, err := gphotos.NewClient(
		gphotos.AuthenticateUser(
			config.OAuthConfig(folderUploadJob.uploaderConfigAPICredentials),
			gphotos.WithUserLoginHint(folderUploadJob.Account),
		),
	)
	if err != nil {
		return nil, errors.Annotate(err, "failed authenticating new client")
	}

	// and store the token into the keyring
	err = tkm.StoreToken(folderUploadJob.Account, gphotosClient.Token())
	if err != nil {
		return nil, errors.Annotate(err, "failed storing token")
	}

	return gphotosClient, nil
}

// Upload uploads folder
func (job *FolderUploadJob) Upload(uploadChan chan<- *FileUpload) error {
	folderAbsolutePath, err := cp.AbsolutePath(job.SourceFolder)
	if err != nil {
		return err
	}

	if !filesystem.IsDir(folderAbsolutePath) {
		return fmt.Errorf("%s is not a folder", folderAbsolutePath)
	}

	// dirs are walked depth-first.   These vars hold the active album
	// default empty album for makeAlbums.enabled = false
	currentPhotoAlbum := &photoslibrary.Album{}
	errW := filepath.Walk(folderAbsolutePath, func(fp string, fi os.FileInfo, errP error) error {
		// log.Printf("Upload.Walk: %v, fi: %v, err: %v\n", fp, fi, err)
		// TODO: integrate error reporting
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error for %v: %v\n", fp, err)
			return nil
		}
		if fi == nil {
			_, _ = fmt.Fprintf(os.Stderr, "error for %v: FileInfo is nil\n", fp)
			return nil
		}

		selectExcludeFilter := func(item string) bool {
			matched, err := filter.List(job.ExcludePatterns, item)
			if err != nil {
				log.Printf("error for exclude pattern: %v", err)
			}

			return !matched
		}

		selectIncludeFilter := func(item string) bool {
			matched, err := filter.List(job.IncludePatterns, item)
			if err != nil {
				log.Printf("error for include pattern: %v", err)
			}

			return matched
		}

		if !selectFunc(selectIncludeFilter, selectExcludeFilter, fp, fi) {
			log.Printf("Upload.Walk: path %v excluded", fp)
			if fi.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}


		if fi.IsDir() {
			if job.MakeAlbums.Enabled && job.MakeAlbums.Use == "folderNames" {
				log.Printf("Entering Directory: %s", fp)
				currentPhotoAlbum, err = job.gphotosClient.GetOrCreateAlbumByName(filepath.Base(fp))
				if err != nil {
					currentPhotoAlbum = &photoslibrary.Album{}
					log.Printf("error creating album: %s. File will be uploaded without album", fp)
					return nil
				}
				log.Printf("using album: %s / Id: %s", currentPhotoAlbum.Title, currentPhotoAlbum.Id)
			} else {
				log.Printf("album not created: %s. set jobs.makeAlbums.enabled = true to create albums", fp)
			}
			return nil
		}
		// process only files
		if !filesystem.IsFile(fp) {
			return nil
		}

		// if we don't upload videos check it's not a video
		if !job.UploadVideos && filetypes.IsVideo(fp) {
			log.Printf("recognized as video file: %s: skipping file... (set uploadVideos to true in config to upload videos)\n", fp)
			return nil
		}

		// check upload db for previous uploads
		isAlreadyUploaded, err := job.completedUploads.IsAlreadyUploaded(fp)
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

		// set file upload options depending on folder upload options
		var fileUpload = &FileUpload{
			FolderUploadJob: job,
			filePath:        fp,
			typedMedia:      typedMedia,
			gphotosClient:   job.gphotosClient.Client,
			album:           currentPhotoAlbum,
		}

		// finally, add the file upload to the queue
		uploadChan <- fileUpload

		return nil
	})
	if errW != nil {
		log.Printf("walk error [%v]\n", errW)
	}

	return nil
}

// selectFunc returns true for all items that should be included (files and
// dirs). If false is returned, files are ignored and dirs are not even walked.
func selectFunc(includeFunc, excludeFunc func(string) bool, item string, fi os.FileInfo) bool {
	return includeFunc(item) && !excludeFunc(item)
}

