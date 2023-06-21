package app

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/spf13/afero"
	"golang.org/x/oauth2"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/config"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/filetracker"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/tokenmanager"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/upload_tracker"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
)

const (
	// DefaultConfigFilename is the default config file name.
	DefaultConfigFilename = "config.hjson"
)

// App represents a running application with all the dependant services.
type App struct {
	// FileTracker tracks local files already uploaded.
	FileTracker FileTracker
	// TokenManager keeps secrets (like tokens).
	TokenManager TokenManager
	// UploadSessionTracker tracks uploads sessions to implement resumable uploads.
	UploadSessionTracker UploadSessionTracker

	// Client is the HTTP client after authentication.
	Client *http.Client

	Logger log.Logger

	// fs points to the file system.
	// Useful for testing.
	fs afero.Fs

	// appDir is the directory application directory.
	appDir string

	// Config keeps the application configuration.
	Config *config.Config
}

// Start initializes the application with the services defined by a given configuration.
// The provided path is the expanded and absolute path to the application data folder.
func Start(ctx context.Context, path string) (*App, error) {
	app, err := StartServices(ctx, path)
	if err != nil {
		return nil, err
	}

	app.Client, err = app.AuthenticateFromToken(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// StartServices initializes the services defined by a given configuration.
// The provided path is the expanded and absolute path to the application data folder.
func StartServices(ctx context.Context, path string) (*App, error) {
	var err error

	app := &App{
		appDir: path,
		Logger: log.GetInstance(),
		fs:     afero.NewOsFs(),
	}

	app.Logger.Infof("Reading configuration from '%s'", app.configFilename())
	app.Config, err = config.FromFile(app.fs, app.configFilename())
	if err != nil {
		return nil, fmt.Errorf("invalid configuration at '%s': %s", app.configFilename(), err)
	}

	app.Logger.Debugf("Current configuration: %s", app.Config.SafePrint())

	if err := app.startServices(); err != nil {
		return nil, err
	}

	return app, nil
}

// StartWithoutConfig initializes the application without reading the configuration.
// The provided path is the expanded and absolute path to the application data folder.
func StartWithoutConfig(fs afero.Fs, path string) (*App, error) {
	app := &App{
		appDir: path,
		Logger: log.GetInstance(),
		fs:     fs,
	}

	app.Logger.Infof("Using application data at '%s'.", app.appDir)

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
	app.UploadSessionTracker.Close()

	// Close token manager
	app.Logger.Debug("Shutting down Token Manager service...")
	if err := app.TokenManager.Close(); err != nil {
		return err
	}

	app.Logger.Debug("All services have been shut down successfully")
	return nil
}

// CreateAppDataDir return the filename after creating the application directory and the configuration file with defaults.
// CreateAppDataDir destroys previous application directory.
func (app *App) CreateAppDataDir() (string, error) {
	if err := app.emptyDir(app.appDir); err != nil {
		return "", err
	}
	filename := app.configFilename()
	_, err := config.Create(app.fs, filename)
	if err != nil {
		return "", err
	}
	return filename, nil
}

// AppDataDirExists return true if the application data dir exists.
func (app *App) AppDataDirExists() bool {
	exist, err := afero.Exists(app.fs, app.configFilename())
	if err != nil {
		return false
	}
	return exist
}

func (app *App) configFilename() string {
	return filepath.Join(app.appDir, DefaultConfigFilename)
}

func (app *App) startServices() error {
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

func (app *App) defaultFileTracker() (*filetracker.FileTracker, error) {
	fileTrackerFolder := filepath.Join(app.appDir, "uploaded_files")
	repo, err := filetracker.NewLevelDBRepository(fileTrackerFolder)
	if err != nil {
		return nil, err
	}
	return filetracker.New(repo), nil
}

func (app *App) defaultTokenManager(backendType string) (*tokenmanager.TokenManager, error) {
	tokensFolder := filepath.Join(app.appDir, "tokens")
	kr, err := tokenmanager.NewKeyringRepository(backendType, nil, tokensFolder)
	if err != nil {
		return nil, err
	}
	return tokenmanager.New(kr), nil
}

func (app *App) defaultUploadsSessionTracker() (*upload_tracker.LevelDBStore, error) {
	ongoingUploadsTrackerFolder := filepath.Join(app.appDir, "ongoing_uploads")
	return upload_tracker.NewStore(ongoingUploadsTrackerFolder)
}

func (app *App) emptyDir(path string) error {
	if err := app.fs.RemoveAll(path); err != nil {
		return err
	}
	return app.fs.MkdirAll(path, 0700)
}

// FileTracker represents a service to track file already uploaded.
type FileTracker interface {
	Put(file string) error
	Exist(file string) bool
	Delete(file string) error
	Close() error
}

// TokenManager represents a service to keep and read secrets (like passwords, tokens...)
type TokenManager interface {
	Put(email string, token *oauth2.Token) error
	Get(email string) (*oauth2.Token, error)
	Close() error
}

// UploadSessionTracker represents a service to keep resumable upload sessions.
//
// See [gphotosuploader/google-photos-api-client-go] Store interface.
type UploadSessionTracker interface {
	Get(fingerprint string) (string, bool)
	Set(fingerprint string, url string)
	Delete(fingerprint string)
	Close()
}
