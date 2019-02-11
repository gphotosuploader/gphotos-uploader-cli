package upload

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"

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

// Upload uploads a folder
func (folderUploadJob *FolderUploadJob) Upload() error {
	folderAbsolutePath, err := cp.AbsolutePath(folderUploadJob.SourceFolder)
	if err != nil {
		return err
	}

	if !filesystem.IsDir(folderAbsolutePath) {
		return fmt.Errorf("%s is not a folder", folderAbsolutePath)
	}

	err = filepath.Walk(folderAbsolutePath, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		// process only files
		if !filesystem.IsFile(filePath) {
			return nil
		}
		// process only media
		if !filetypes.IsMedia(filePath) {
			fmt.Printf("not a media file: %s: skipping file...\n", filePath)
			return nil
		}

		// if we don't upload videos check it's not a video
		if !folderUploadJob.UploadVideos && filetypes.IsVideo(filePath) {
			fmt.Printf("recognized as video file: %s: skipping file... (set uploadVideos to true in config to upload videos)\n", filePath)
			return nil
		}

		// check upload db for previous uploads
		isAlreadyUploaded, err := folderUploadJob.completedUploads.IsAlreadyUploaded(filePath)
		if err != nil {
			log.Println(err)
		} else if isAlreadyUploaded {
			fmt.Printf("already uploaded: %s: skipping file...\n", filePath)
			return nil
		}

		typedMedia, err := filetypes.NewTypedMedia(filePath)
		if err != nil {
			fmt.Println(errors.Annotatef(err, "failed creating new TypedMedia from filePath"))
			return nil
		}

		// set file upload options depending on folder upload options
		var fileUpload = &FileUpload{
			FolderUploadJob: folderUploadJob,
			filePath:        filePath,
			typedMedia:      typedMedia,
			gphotosClient:   folderUploadJob.gphotosClient.Client,
		}
		if folderUploadJob.MakeAlbums.Enabled && folderUploadJob.MakeAlbums.Use == USEFOLDERNAMES {
			lastDirName := filepath.Base(filepath.Dir(filePath))
			fileUpload.albumName = lastDirName
		}

		// finally, add the file upload to the queue
		QueueFileUpload(fileUpload)

		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}

	return nil
}
