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

const (
	// DefaultConfigFilename is the default config File name to use
	DefaultConfigFilename = "config.hjson"
)

// Os points to the (real) file system.
// Useful for testing.
var Os = afero.NewOsFs()

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

// Create returns the configuration file name created with default settings.
func Create(dir string) (string, error) {
	cfg := defaultSettings()
	file := defaultConfigFilePath(dir)
	if err := cfg.writeFile(file); err != nil {
		return "", err
	}
	return file, nil
}

// FromFile returns the configuration data read from the specified directory.
// FromFile returns a ParseError{} if the configuration validation fails.
func FromFile(dir string) (*Config, error) {
	filename := defaultConfigFilePath(dir)
	cfg, err := readFile(filename)
	if err != nil {
		return nil, ParseError{err}
	}
	if err := cfg.validate(); err != nil {
		return cfg, ParseError{err}
	}

	return cfg, nil
}

// Exists checks the existence of the configuration file
func Exists(path string) bool {
	file := defaultConfigFilePath(path)
	path = normalizePath(path)
	if _, err := Os.Stat(file); err != nil {
		return false
	}
	return true
}

// validate returns if the current configuration is valid.
func (c *Config) validate() error {
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

// writeFile writes the configuration data to a file named by filename.
// If the file does not exist, writeFile creates it;
// otherwise writeFile truncates it before writing.
func (c Config) writeFile(filename string) error {
	b, err := hjson.MarshalWithOptions(c, hjson.DefaultOptions())
	if err != nil {
		return err

	}
	return afero.WriteFile(Os, filename, b, 0600)
}

// readFile loads the configuration data reading the file named by filename.
func readFile(filename string) (*Config, error) {
	b, err := afero.ReadFile(Os, filename)
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
		return errors.New("config: APIAppCredentials are invalid")
	}
	return nil
}

func (c Config) validateAccount() error {
	if c.Account == "" {
		return errors.New("config: Account could not be empty")
	}
	return nil
}

func (c Config) validateJobs() error {
	if len(c.Jobs) < 1 {
		return errors.New("config: At least one Job must be configured")
	}

	for _, job := range c.Jobs {
		exist, err := afero.DirExists(Os, job.SourceFolder)
		if err != nil {
			return fmt.Errorf("config: The provided folder '%s' could not be used, err=%s", job.SourceFolder, err)
		}
		if !exist {
			return fmt.Errorf("config: The provided folder '%s' is not a folder", job.SourceFolder)
		}
		if job.MakeAlbums.Enabled &&
			(job.MakeAlbums.Use != "folderPath" && job.MakeAlbums.Use != "folderName") {
			return errors.New("config: The provided MakeAlbums option is invalid")
		}
	}
	return nil
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

func defaultConfigFilePath(path string) string {
	path = filepath.Join(path, DefaultConfigFilename)
	return normalizePath(path)
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

func (c Config) validateSecretsBackendType() error {
	if c.SecretsBackendType != "auto" &&
		c.SecretsBackendType != "secret-service" &&
		c.SecretsBackendType != "keychain" &&
		c.SecretsBackendType != "kwallet" &&
		c.SecretsBackendType != "file" {
		return fmt.Errorf("config: SecretsBackendType is invalid, %s", c.SecretsBackendType)
	}
	return nil
}

func normalizePath(path string) string {
	if absPath, err := filepath.Abs(path); err == nil {
		return absPath
	}
	return filepath.Clean(path)
}
