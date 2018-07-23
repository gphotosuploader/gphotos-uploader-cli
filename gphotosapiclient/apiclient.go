package gphotosapiclient

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/palantir/stacktrace"

	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

const apiVersion = "v1"
const basePath = "https://photoslibrary.googleapis.com/"

// PhotosClient is a client for uploading a media.
// photoslibrary does not provide `/v1/uploads` API so we implement here.
type PhotosClient struct {
	*photoslibrary.Service
	Client *http.Client
}

// TODO: FIGURE OUT / REMOVE
// //////////////////////////
type PhotosClientConstructorOptionsContainer struct {
	FromKeyring bool
}
type PhotosClientConstructorOption func(options *PhotosClientConstructorOptionsContainer) *PhotosClientConstructorOptionsContainer

func FromKeyring(options *PhotosClientConstructorOptionsContainer) *PhotosClientConstructorOptionsContainer {
	options.FromKeyring = true
	return options
}

// /////////////////////////////

// New constructs a new PhotosClient from an httpClient
func New(httpClient *http.Client, options ...PhotosClientConstructorOption) (photosClient *PhotosClient, err error) {
	photosLibraryClient, err := photoslibrary.New(httpClient)
	if err != nil {
		return nil, err
	}
	return &PhotosClient{photosLibraryClient, httpClient}, nil
}

// GetUploadToken sends the media and returns the UploadToken.
func (client *PhotosClient) GetUploadToken(r io.Reader, filename string) (token string, err error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/uploads", basePath, apiVersion), r)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-Goog-Upload-File-Name", filename)

	res, err := client.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	uploadToken := string(b)
	return uploadToken, nil
}

// Upload actually uploads the media and activates it on google photos
func (client *PhotosClient) Upload(filePath string) error {
	filename := path.Base(filePath)
	log.Printf("Uploading %s", filename)

	file, err := os.Open(filePath)
	if err != nil {
		return stacktrace.Propagate(err, "failed opening file")
	}
	defer file.Close()

	uploadToken, err := client.GetUploadToken(file, filename)
	if err != nil {
		return stacktrace.Propagate(err, "failed getting uploadToken for %s", filename)
	}
	log.Printf("Uploaded %s as token %s", filename, uploadToken)

	log.Printf("Adding media %s", filename)
	batch, err := client.MediaItems.BatchCreate(&photoslibrary.BatchCreateMediaItemsRequest{
		NewMediaItems: []*photoslibrary.NewMediaItem{
			&photoslibrary.NewMediaItem{
				Description:     filename,
				SimpleMediaItem: &photoslibrary.SimpleMediaItem{UploadToken: uploadToken},
			},
		},
	}).Do()
	if err != nil {
		return stacktrace.Propagate(err, "failed adding media %s", filename)
	}

	for _, result := range batch.NewMediaItemResults {
		log.Printf("Added media %s as %s", result.MediaItem.Description, result.MediaItem.Id)
	}

	return nil
}
