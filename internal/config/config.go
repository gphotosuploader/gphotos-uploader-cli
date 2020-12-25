package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"

	"github.com/hjson/hjson-go"
	"github.com/spf13/afero"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/filesystem"
)

const (
	// DefaultConfigFilename is the default config File name to use
	DefaultConfigFilename = "config.hjson"
)

// ParseError denotes failing to parse configuration file.
type ParseError struct {
	err error
}

// Error returns the formatted configuration error.
func (e ParseError) Error() string {
	return fmt.Sprintf("parsing config: %s", e.err.Error())
}

// CreateError denotes failing to create configuration file.
type CreateError struct {
	path string
	err  error
}

// Error returns the formatted configuration error.
func (e CreateError) Error() string {
	return fmt.Sprintf("creating config path '%s': %s", e.path, e.err.Error())

}

// CompletedUploadsDBDir returns the path of the folder where completed uploads are tracked.
func (c *AppConfig) CompletedUploadsDBDir() string {
	return path.Join(c.ConfigPath, "uploads.db")
}

// ResumableUploadsDBDir returns the path of the folder where upload URLs are tracked.
func (c *AppConfig) ResumableUploadsDBDir() string {
	return path.Join(c.ConfigPath, "resumable_uploads.db")
}

// File return the path of the configuration File.
func (c *AppConfig) File() string {
	return path.Join(c.ConfigPath, DefaultConfigFilename)
}

// KeyringDir returns the path of the folder where keyring will be stored.
// This is only used if 'SecretsBackendType=File'
func (c *AppConfig) KeyringDir() string {
	return c.ConfigPath
}

// New creates a new configuration file path dir with default settings.
// New removes all configuration inside the specified dir.
func New(fs afero.Fs, dir string) (*AppConfig, error) {
	configDefaults := defaultSettings()
	cfg := AppConfig{
		ConfigPath: ensureAbsolutePath(dir),
		Config:     &configDefaults,
	}

	// Empty the application directory.
	cfgPath := ensureAbsolutePath(dir)
	if err := filesystem.EmptyOrCreateDir(cfgPath); err != nil {
		return nil, CreateError{
			path: cfgPath,
			err:  err,
		}
	}

	if err := cfg.writeFile(fs, filepath.Join(cfgPath, DefaultConfigFilename)); err != nil {
		return nil, CreateError{
			path: cfgPath,
			err:  err,
		}
	}

	return &cfg, nil
}

// FromFile returns the configuration data read from the specified directory.
// FromFile returns a ParseError{} if the configuration validation fails.
func FromFile(fs afero.Fs, dir string) (*AppConfig, error) {
	cfg := &AppConfig{
		ConfigPath: ensureAbsolutePath(dir),
	}

	filename := filepath.Join(ensureAbsolutePath(dir), DefaultConfigFilename)
	if err := cfg.readFile(fs, filename); err != nil {
		return cfg, ParseError{err}
	}
	if err := cfg.Validate(); err != nil {
		return cfg, ParseError{err}
	}
	return cfg, nil
}

// Validate returns if the current configuration is valid.
func (c *AppConfig) Validate() error {
	if err := c.validateSecretsBackendType(); err != nil {
		return err
	}
	if err := c.validateAPIAppCredentials(); err != nil {
		return err
	}
	if err := c.validateAccount(); err != nil {
		return err
	}
	if err := c.validateJobs(); err != nil {
		return err
	}
	return nil
}

// Exists checks if a gphotos-uplaoder-cli configuration exists at a certain path
func Exists(fs afero.Fs, path string) bool {
	cfgFile, err := filesystem.AbsolutePath(filepath.Join(path, DefaultConfigFilename))
	if err != nil {
		return false
	}

	if _, err := fs.Stat(cfgFile); err == nil {
		return true
	}

	return false
}

// writeFile writes the configuration data to a file named by filename.
// If the file does not exist, writeFile creates it;
// otherwise writeFile truncates it before writing.
func (c *AppConfig) writeFile(fs afero.Fs, filename string) error {
	b, err := hjson.MarshalWithOptions(c.Config, hjson.DefaultOptions())
	if err != nil {
		return err

	}
	return afero.WriteFile(fs, filename, b, 0600)
}

// readFile loads the configuration data reading the file named by filename.
func (c *AppConfig) readFile(fs afero.Fs, filename string) error {
	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		return err
	}

	if err := unmarshalReader(bytes.NewReader(b), c.Config); err != nil {
		return err
	}

	// convert all path to absolute paths.
	return c.ensureSourceFolderAbsolutePaths()
}

func (c *AppConfig) validateAPIAppCredentials() error {
	if c.APIAppCredentials.ClientID == "" || c.APIAppCredentials.ClientSecret == "" {
		return errors.New("config: APIAppCredentials are invalid")
	}
	return nil
}

func (c *AppConfig) validateAccount() error {
	if c.Account == "" {
		return errors.New("config: Account could not be empty")
	}
	return nil
}

func (c *AppConfig) validateJobs() error {
	if len(c.Jobs) < 1 {
		return errors.New("config: At least one Job must be configured")
	}

	for _, job := range c.Jobs {
		if !filesystem.IsDir(job.SourceFolder) {
			return fmt.Errorf("config: The provided SourceFolder is not a folder. [%s]", job.SourceFolder)
		}
		if job.MakeAlbums.Enabled &&
			(job.MakeAlbums.Use != "folderPath" && job.MakeAlbums.Use != "folderName") {
			return errors.New("config: The provided MakeAlbums option is invalid")
		}
	}
	return nil
}

func (c *AppConfig) ensureSourceFolderAbsolutePaths() error {
	for i := range c.Jobs {
		item := &c.Jobs[i] // we do that way to modify original object while iterating.
		srcFolder, err := filesystem.AbsolutePath(item.SourceFolder)
		if err != nil {
			return ParseError{err}
		}
		item.SourceFolder = srcFolder
	}
	return nil
}

// unmarshalReader unmarshal HJSON data.
func unmarshalReader(in io.Reader, c interface{}) error {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(in)

	b, err := hjsonToJson(buf.Bytes())
	if err != nil {
		return err
	}

	// unmarshal
	return json.Unmarshal(b, c)
}

// hjsonToJson converts dta from HJSON to JSON format.
func hjsonToJson(in []byte) ([]byte, error) {
	var raw map[string]interface{}
	if err := hjson.Unmarshal(in, &raw); err != nil {
		return nil, err
	}

	// convert to JSON
	return json.Marshal(raw)
}

// defaultSettings() returns a *AppConfig with the default settings of the application.
func defaultSettings() Config {
	return Config{
		SecretsBackendType: "auto",
		APIAppCredentials: APIAppCredentials{
			ClientID:     "20637643488-1hvg8ev08r4tc16ca7j9oj3686lcf0el.apps.googleusercontent.com",
			ClientSecret: "0JyfLYw0kyDcJO-pGg5-rW_P",
		},
		Account: "youremail@gmail.com",
		Jobs: []FolderUploadJob{
			{
				SourceFolder: "~/folder/to/upload",
				MakeAlbums: MakeAlbums{
					Enabled: true,
					Use:     "folderName",
				},
				DeleteAfterUpload: false,
			},
		},
	}
}

func (c *AppConfig) validateSecretsBackendType() error {
	if c.SecretsBackendType != "auto" &&
		c.SecretsBackendType != "secret-service" &&
		c.SecretsBackendType != "keychain" &&
		c.SecretsBackendType != "kwallet" &&
		c.SecretsBackendType != "file" {
		return fmt.Errorf("config: SecretsBackendType is invalid, %s", c.SecretsBackendType)
	}
	return nil
}

func ensureAbsolutePath(path string) string {
	absPath, err := filesystem.AbsolutePath(path)
	if err != nil {
		return path
	}
	return absPath
}
