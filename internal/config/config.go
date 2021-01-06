package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/hjson/hjson-go"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

// ParseError denotes failing to parse configuration file.
type ParseError struct {
	err error
}

// Error returns the formatted configuration error.
func (e ParseError) Error() string {
	return fmt.Sprintf("parsing config: %s", e.err.Error())
}

// Create returns the configuration data after creating file with default settings.
func Create(fs afero.Fs, filename string) (*Config, error) {
	cfg := defaultSettings()
	if err := cfg.writeFile(fs, filename); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// FromFile returns the configuration data read from the specified file.
// FromFile returns a ParseError{} if the configuration validation fails.
func FromFile(fs afero.Fs, filename string) (*Config, error) {
	cfg, err := readFile(fs, filename)
	if err != nil {
		return nil, ParseError{err}
	}
	if err := cfg.validate(fs); err != nil {
		return cfg, ParseError{err}
	}

	return cfg, nil
}

// Exists checks the existence of the configuration file
func Exists(fs afero.Fs, filename string) bool {
	filename = normalizePath(filename)
	if _, err := fs.Stat(filename); err != nil {
		return false
	}
	return true
}

// validate validates the current configuration.
func (c Config) validate(fs afero.Fs) error {
	if err := c.validateSecretsBackendType(); err != nil {
		return err
	}
	if err := c.validateAPIAppCredentials(); err != nil {
		return err
	}
	if err := c.validateAccount(); err != nil {
		return err
	}
	if err := c.validateJobs(fs); err != nil {
		return err
	}
	return nil
}

// writeFile writes the configuration data to a file named by filename.
// If the file does not exist, writeFile creates it;
// otherwise writeFile truncates it before writing.
func (c Config) writeFile(fs afero.Fs, filename string) error {
	b, err := hjson.MarshalWithOptions(c, hjson.DefaultOptions())
	if err != nil {
		return err

	}
	return afero.WriteFile(fs, filename, b, 0600)
}

// readFile loads the configuration data reading the file named by filename.
func readFile(fs afero.Fs, filename string) (*Config, error) {
	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := unmarshalReader(bytes.NewReader(b), &config); err != nil {
		return nil, err
	}

	// convert all path to absolute paths.
	if err := config.ensureSourceFolderAbsolutePaths(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c Config) validateAPIAppCredentials() error {
	if c.APIAppCredentials.ClientID == "" || c.APIAppCredentials.ClientSecret == "" {
		return errors.New("option APIAppCredentials are invalid")
	}
	return nil
}

func (c Config) validateAccount() error {
	if c.Account == "" {
		return errors.New("option Account could not be empty")
	}
	return nil
}

func (c Config) validateJobs(fs afero.Fs) error {
	if len(c.Jobs) < 1 {
		return errors.New("at least one Job must be configured")
	}

	for _, job := range c.Jobs {
		exist, err := afero.DirExists(fs, job.SourceFolder)
		if err != nil {
			return fmt.Errorf("option SourceFolder '%s' is invalid, err=%s", job.SourceFolder, err)
		}
		if !exist {
			return fmt.Errorf("option SourceFolder '%s' does not exist", job.SourceFolder)
		}
		if job.MakeAlbums.Enabled &&
			(job.MakeAlbums.Use != "folderPath" && job.MakeAlbums.Use != "folderName") {
			return fmt.Errorf("option MakeAlbums is invalid, '%s", job.MakeAlbums.Use)
		}
	}
	return nil
}

func (c Config) validateSecretsBackendType() error {
	switch c.SecretsBackendType {
	case "auto", "secret-service", "keychain", "kwallet", "file":
		return nil

	}
	return fmt.Errorf("option SecretsBackendType is invalid, '%s'", c.SecretsBackendType)
}

func (c Config) ensureSourceFolderAbsolutePaths() error {
	for i := range c.Jobs {
		item := &c.Jobs[i] // we do that way to modify original object while iterating.
		src, err := homedir.Expand(item.SourceFolder)
		if err != nil {
			return ParseError{err}
		}
		item.SourceFolder = normalizePath(src)
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

// defaultSettings() returns a *Config with the default settings of the application.
func defaultSettings() Config {
	return Config{
		SecretsBackendType: "file",
		APIAppCredentials: APIAppCredentials{
			ClientID:     "YOUR_APP_CLIENT_ID",
			ClientSecret: "YOUR_APP_CLIENT_SECRET",
		},
		Account: "YOUR_GOOGLE_PHOTOS_ACCOUNT",
		Jobs: []FolderUploadJob{
			{
				SourceFolder: "YOUR_FOLDER_PATH",
				MakeAlbums: MakeAlbums{
					Enabled: true,
					Use:     "folderName",
				},
				DeleteAfterUpload: false,
			},
		},
	}
}

func normalizePath(path string) string {
	if absPath, err := filepath.Abs(path); err == nil {
		return absPath
	}
	return filepath.Clean(path)
}
