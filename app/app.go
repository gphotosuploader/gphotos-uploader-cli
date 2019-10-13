package app

import (
	"log"

	"golang.org/x/oauth2"
)

type App struct {
	FileTracker   FileTracker
	TokenManager  TokenManager
	UploadTracker UploadTracker

	Log *log.Logger
}

type FileTracker interface {
	CacheAsAlreadyUploaded(filePath string) error
	IsAlreadyUploaded(filePath string) (bool, error)
	RemoveAsAlreadyUploaded(filePath string) error
}

type TokenManager interface {
	StoreToken(email string, token *oauth2.Token) error
	RetrieveToken(email string) (*oauth2.Token, error)
}

type UploadTracker interface {
	Get(fingerprint string) []byte
	Set(fingerprint string, url []byte)
	Delete(fingerprint string)
}
