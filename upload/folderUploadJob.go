package upload

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/palantir/stacktrace"

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

func (folderUploadJob *FolderUploadJob) Run() {
	sourceFolderAbsolutePath, err := cp.AbsolutePath(folderUploadJob.SourceFolder)
	if err != nil {
		log.Fatal(err)
	}

	client, err := Authenticate(folderUploadJob)
	if err != nil {
		log.Fatal(err)
	}

	err = folderUploadJob.uploadFolder(client, sourceFolderAbsolutePath)
	if err != nil {
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
	gphotosClient, err := gphotos.NewClient(
		gphotos.AuthenticateUser(config.OAuthConfig()),
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

func (j *FolderUploadJob) uploadFolder(gphotosClient *gphotos.Client, folderPath string) error {
	if !fileshandling.IsDir(folderPath) {
		return fmt.Errorf("%s is not a folder", folderPath)
	}

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if fileshandling.IsFile(path) {
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
