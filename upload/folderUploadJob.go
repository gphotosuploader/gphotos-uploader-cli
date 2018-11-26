package upload

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/palantir/stacktrace"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/nmrshll/go-cp"
	gphotos "github.com/nmrshll/google-photos-api-client-go/noserver-gphotos"
	"github.com/nmrshll/gphotos-uploader-cli/config"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/tokenstore"
	"github.com/nmrshll/gphotos-uploader-cli/fileshandling"
)

const (
	USEFOLDERNAMES = "folderNames"
)

type FolderUploadJob struct {
	*config.FolderUploadJob
}

func (folderUploadJob *FolderUploadJob) Run(db *leveldb.DB) {
	sourceFolderAbsolutePath, err := cp.AbsolutePath(folderUploadJob.SourceFolder)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	client, err := Authenticate(folderUploadJob)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	err = folderUploadJob.uploadFolder(client, sourceFolderAbsolutePath, db)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
}

func Authenticate(folderUploadJob *FolderUploadJob) (*gphotos.Client, error) {
	// try to load token from keyring
	token, err := tokenstore.RetrieveToken(folderUploadJob.Account)
	if err == nil && token != nil { // if error ignore and skip
		// if found create client from token
		gphotosClient, err := gphotos.NewClient(gphotos.FromToken(config.OAuthConfig(), token))
		if err == nil && gphotosClient != nil { // if error ignore and skip
			return gphotosClient, nil
		}
	}

	// else authenticate again to grab a new token
	log.Println(color.CyanString(fmt.Sprintf("Need to log login into account %s", folderUploadJob.Account)))
	time.Sleep(1200 * time.Millisecond)
	gphotosClient, err := gphotos.NewClient(
		gphotos.AuthenticateUser(
			config.OAuthConfig(),
			gphotos.WithUserLoginHint(folderUploadJob.Account),
		),
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed authenticating new client")
	}

	// and store the token into the keyring
	err = tokenstore.StoreToken(folderUploadJob.Account, gphotosClient.Token())
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed storing token")
	}

	return gphotosClient, nil
}

func (j *FolderUploadJob) uploadFolder(gphotosClient *gphotos.Client, folderPath string, db *leveldb.DB) error {
	if !fileshandling.IsDir(folderPath) {
		return fmt.Errorf("%s is not a folder", folderPath)
	}

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if fileshandling.IsFile(path) {
			// check type is image
			if isImage, err := fileshandling.IsImage(path); err != nil || !isImage {
				fmt.Printf("not an image: %s: skipping file...\n", path)
				return nil
			}

			// check upload db for previous uploads
			isUploaded, err := fileshandling.IsUploadedPrev(path, db)
			if err != nil {
				log.Println(err)
			} else if isUploaded {
				fmt.Printf("previously uploaded: %s: skipping file...\n", path)
				return nil
			}

			var fileUpload = &FileUpload{FolderUploadJob: j, filePath: path, gphotosClient: gphotosClient.Client}
			if j.MakeAlbums.Enabled && j.MakeAlbums.Use == USEFOLDERNAMES {
				lastDirName := filepath.Base(filepath.Dir(path))
				fileUpload.albumName = lastDirName
			}
			QueueFileUpload(fileUpload)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	return nil
}
