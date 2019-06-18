package upload

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	photoslibrary "google.golang.org/api/photoslibrary/v1"

	"github.com/juju/errors"
	cp "github.com/nmrshll/go-cp"
	gphotos "github.com/nmrshll/google-photos-api-client-go/noserver-gphotos"
	"github.com/nmrshll/gphotos-uploader-cli/config"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/completeduploads"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/tokenstore"
	"github.com/nmrshll/gphotos-uploader-cli/filetypes"
	"github.com/nmrshll/gphotos-uploader-cli/utils/filesystem"
)

const (
	USEFOLDERNAMES = "folderNames"
)

type FolderUploadJob struct {
	*config.FolderUploadJob
	uploaderConfigAPICredentials *config.APIAppCredentials
	gphotosClient                *gphotos.Client
	completedUploads             *completeduploads.CompletedUploadsService
}

func NewFolderUploadJob(configFolderUploadJob *config.FolderUploadJob, completedUploads *completeduploads.CompletedUploadsService, uploaderConfigAPICredentials *config.APIAppCredentials) *FolderUploadJob {
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

	gphotosClient, err := authenticate(folderUploadJob)
	if err != nil {
		log.Fatal(err)
	}
	folderUploadJob.gphotosClient = gphotosClient

	return folderUploadJob
}

func authenticate(folderUploadJob *FolderUploadJob) (*gphotos.Client, error) {
	// try to load token from keyring
	token, err := tokenstore.RetrieveToken(folderUploadJob.Account)
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
	err = tokenstore.StoreToken(folderUploadJob.Account, gphotosClient.Token())
	if err != nil {
		return nil, errors.Annotate(err, "failed storing token")
	}

	return gphotosClient, nil
}

// Upload uploads folder
func (folderUploadJob *FolderUploadJob) Upload(fileUploadsChan chan<- *FileUpload) error {
	folderAbsolutePath, err := cp.AbsolutePath(folderUploadJob.SourceFolder)
	if err != nil {
		return err
	}

	if !filesystem.IsDir(folderAbsolutePath) {
		return fmt.Errorf("%s is not a folder", folderAbsolutePath)
	}

	// dirs are walked depth-first.   These vars hold the active album
	// default empty album for makeAlbums.enabled = false
	currentPhotoAlbum := &photoslibrary.Album{}
	errW := filepath.Walk(folderAbsolutePath, func(filePath string, info os.FileInfo, errP error) error {
		if info.IsDir() {
			if folderUploadJob.MakeAlbums.Enabled && folderUploadJob.MakeAlbums.Use == USEFOLDERNAMES {
				log.Printf("Entering Directory: %s", filePath)
				currentPhotoAlbum, err = folderUploadJob.gphotosClient.GetOrCreateAlbumByName(filepath.Base(filePath))
				if err != nil {
					currentPhotoAlbum = &photoslibrary.Album{}
					log.Printf("error creating album: %s. File will be uploaded without album", filePath)
					return nil
				}
				log.Printf("using album: %s / Id: %s", currentPhotoAlbum.Title, currentPhotoAlbum.Id)
			} else {
				log.Printf("album not created: %s. set jobs.makeAlbums.enabled = true to create albums", filePath)
			}
			return nil
		}
		// process only files
		if !filesystem.IsFile(filePath) {
			return nil
		}
		// process only media
		if !filetypes.IsMedia(filePath) {
			log.Printf("not a media file: %s: skipping file...\n", filePath)
			return nil
		}

		// if we don't upload videos check it's not a video
		if !folderUploadJob.UploadVideos && filetypes.IsVideo(filePath) {
			log.Printf("recognized as video file: %s: skipping file... (set uploadVideos to true in config to upload videos)\n", filePath)
			return nil
		}

		// check upload db for previous uploads
		isAlreadyUploaded, err := folderUploadJob.completedUploads.IsAlreadyUploaded(filePath)
		if err != nil {
			log.Println(err)
		} else if isAlreadyUploaded {
			log.Printf("already uploaded: %s: skipping file...\n", filePath)
			return nil
		}

		typedMedia, err := filetypes.NewTypedMedia(filePath)
		if err != nil {
			log.Println(errors.Annotatef(err, "failed creating new TypedMedia from filePath"))
			return nil
		}

		// set file upload options depending on folder upload options
		var fileUpload = &FileUpload{
			FolderUploadJob: folderUploadJob,
			filePath:        filePath,
			typedMedia:      typedMedia,
			gphotosClient:   folderUploadJob.gphotosClient.Client,
			album:           currentPhotoAlbum,
		}

		// finally, add the file upload to the queue
		fileUploadsChan <- fileUpload

		return nil
	})
	if errW != nil {
		log.Printf("walk error [%v]\n", errW)
	}

	return nil
}
