package upload

import (
	"fmt"
	"github.com/nmrshll/go-cp"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/noserver-gphotos"
	"github.com/juju/errors"
	"github.com/nmrshll/gphotos-uploader-cli/config"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/completeduploads"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/tokenstore"
	"github.com/nmrshll/gphotos-uploader-cli/filetypes"
	"github.com/nmrshll/gphotos-uploader-cli/utils/filesystem"
)

// Job represents a job to upload all photos in a folder
type Job struct {
	*config.FolderUploadJob
	uploaderConfigAPICredentials *config.APIAppCredentials
	gphotosClient                *gphotos.Client
	trackingRepository           *completeduploads.Service
}

// TODO: accept a *gphotos.Client instead of creating it inside. We can remove a lot of parameters on call and in Job
// NewFolderUploadJob creates a Job based on the submitted data
func NewFolderUploadJob(configFolderUploadJob *config.FolderUploadJob, completedUploads *completeduploads.Service, uploaderConfigAPICredentials *config.APIAppCredentials, tokenManagerService *tokenstore.Service) *Job {
	// check args
	{
		if completedUploads == nil {
			log.Fatalf("completedUploadsService can't be nil")
		}
		if uploaderConfigAPICredentials == nil {
			log.Fatalf("uploaderConfigAPICredentials can't be nil")
		}
	}

	folderUploadJob := &Job{
		FolderUploadJob:              configFolderUploadJob,
		trackingRepository:           completedUploads,
		uploaderConfigAPICredentials: uploaderConfigAPICredentials,
	}

	gphotosClient, err := authenticate(tokenManagerService, folderUploadJob)
	if err != nil {
		log.Fatal(err)
	}
	folderUploadJob.gphotosClient = gphotosClient

	return folderUploadJob
}

// TODO: Move this to a new package Photos or GPhotos where all the Google Photos code is there
func authenticate(tkm *tokenstore.Service, folderUploadJob *Job) (*gphotos.Client, error) {
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
		isAlreadyUploaded, err := job.trackingRepository.IsAlreadyUploaded(fp)
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
			path:            fp,
			typedMedia:      typedMedia,
			gphotosClient:   job.gphotosClient.Client,
			album:           album,
			deleteOnSuccess: job.DeleteAfterUpload,
		}

		// finally, add the file upload to the queue
		uploadChan <- uploadItem

		return nil
	})
	if errW != nil {
		log.Printf("walk error [%v]\n", errW)
	}

	return nil
}
