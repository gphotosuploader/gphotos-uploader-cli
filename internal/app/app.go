package app

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/oauth2"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/config"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/completeduploads"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/leveldbstore"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/tokenstore"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
)

// App represents a running application with all the dependant services.
type App struct {
	FileTracker   FileTracker
	TokenManager  TokenManager
	UploadTracker UploadTracker

	Logger log.Logger
}

// Start initializes the application with the services defined by a given configuration.
func Start(cfg *config.AppConfig) (*App, error) {
	app := &App{}

	// Initialize the logger
	app.Logger = log.GetInstance()

	// Use LevelDB to track already uploaded files
	ft, err := leveldb.OpenFile(cfg.CompletedUploadsDBDir(), nil)
	if err != nil {
		return app, fmt.Errorf("open completed uploads tracker failed: path=%s, err=%s", cfg.CompletedUploadsDBDir(), err)
	}
	app.FileTracker = completeduploads.NewService(completeduploads.NewLevelDBRepository(ft))

	// Use Keyring to store / read secrets
	kr, err := tokenstore.NewKeyringRepository(cfg.SecretsBackendType, nil, cfg.KeyringDir())
	if err != nil {
		return app, fmt.Errorf("open token manager failed: type=%s, err=%s", cfg.SecretsBackendType, err)
	}
	app.TokenManager = tokenstore.New(kr)

	// Upload session tracker to keep upload session to resume uploads.
	app.UploadTracker, err = leveldbstore.NewStore(cfg.ResumableUploadsDBDir())
	if err != nil {
		return app, fmt.Errorf("open resumable uploads tracker failed: path=%s, err=%s", cfg.ResumableUploadsDBDir(), err)
	}

	return app, nil
}

// Stop stops the application releasing all service resources.
func (app *App) Stop() error {
	// Close already uploaded file tracker
	app.Logger.Debug("Shutting down File Tracker service...")
	if err := app.FileTracker.Close(); err != nil {
		return err
	}

	// Close upload session tracker
	app.Logger.Debug("Shutting down Upload Tracker service...")
	if err := app.UploadTracker.Close(); err != nil {
		return err
	}

	// Close token manager
	app.Logger.Debug("Shutting down Token Manager service...")
	if err := app.TokenManager.Close(); err != nil {
		return err
	}

	app.Logger.Debug("All services has been shut down successfully")
	return nil
}

// FileTracker represents a service to track file already uploaded.
type FileTracker interface {
	CacheAsAlreadyUploaded(filePath string) error
	IsAlreadyUploaded(filePath string) (bool, error)
	RemoveAsAlreadyUploaded(filePath string) error
	Close() error
}

// TokenManager represents a service to keep and read secrets (like passwords, tokens...)
type TokenManager interface {
	StoreToken(email string, token *oauth2.Token) error
	RetrieveToken(email string) (*oauth2.Token, error)
	Close() error
}

// UploadTracker represents a service to keep resumable upload sessions.
type UploadTracker interface {
	Get(fingerprint string) []byte
	Set(fingerprint string, url []byte)
	Delete(fingerprint string)
	Close() error
}
