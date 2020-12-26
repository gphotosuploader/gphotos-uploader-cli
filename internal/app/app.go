package app

import (
	"fmt"
	"path/filepath"

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
	// FileTracker tracks local files already uploaded.
	FileTracker FileTracker
	// TokenManager keeps secrets (like tokens).
	TokenManager TokenManager
	// UploadSessionTracker tracks uploads sessions to implement resumable uploads.
	UploadSessionTracker UploadSessionTracker

	Logger log.Logger

	// appDir is the directory application directory.
	appDir string

	// Config keeps the application configuration.
	Config *config.Config
}

// Start initializes the application with the services defined by a given configuration.
func Start(path string) (*App, error) {
	var err error

	app := &App{
		appDir: path,
		Logger: log.GetInstance(),
	}

	app.Logger.Debugf("Reading configuration from '%s'", app.appDir)
	app.Config, err = config.FromFile(app.appDir)
	if err != nil {
		return nil, fmt.Errorf("please review your configuration: file=%s, err=%s", app.appDir, err)
	}

	if err := app.startServices(); err != nil {
		return nil, err
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
	if err := app.UploadSessionTracker.Close(); err != nil {
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

func (app App) startServices() error {
	var err error
	app.FileTracker, err = app.defaultFileTracker()
	if err != nil {
		app.Logger.Errorf("File tracker could not be started, err: %s", err)
		return fmt.Errorf("file tracker could not be started, err: %s", err)
	}
	app.TokenManager, err = app.defaultTokenManager(app.Config.SecretsBackendType)
	if err != nil {
		app.Logger.Errorf("Token manager could not be started, err: %s", err)
		return fmt.Errorf("token manager could not be started, type:%s, err: %s", app.Config.SecretsBackendType, err)
	}
	app.UploadSessionTracker, err = app.defaultUploadsSessionTracker()
	if err != nil {
		app.Logger.Errorf("Uploads session tracker could not be started, err: %s", err)
		return fmt.Errorf("uploads session tracker could not be started, err:%s", err)
	}
	return nil
}

func (app App) defaultFileTracker() (*completeduploads.Service, error) {
	ft, err := leveldb.OpenFile(filepath.Join(app.appDir, "uploads.db"), nil)
	if err != nil {
		return nil, err
	}
	return completeduploads.NewService(completeduploads.NewLevelDBRepository(ft)), nil
}

func (app App) defaultTokenManager(backendType string) (*tokenstore.TokenManager, error) {
	kr, err := tokenstore.NewKeyringRepository(backendType, nil, app.appDir)
	if err != nil {
		return nil, err
	}
	return tokenstore.New(kr), nil
}

func (app App) defaultUploadsSessionTracker() (*leveldbstore.LevelDBStore, error) {
	return leveldbstore.NewStore(filepath.Join(app.appDir, "resumable_uploads.db"))
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

// UploadSessionTracker represents a service to keep resumable upload sessions.
type UploadSessionTracker interface {
	Get(fingerprint string) []byte
	Set(fingerprint string, url []byte)
	Delete(fingerprint string)
	Close() error
}
