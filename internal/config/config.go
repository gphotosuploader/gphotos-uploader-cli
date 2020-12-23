package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/hjson/hjson-go"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/utils/filesystem"
)

const (
	// DefaultConfigFilename is the default config File name to use
	DefaultConfigFilename = "config.hjson"
)

// ParseError denotes failing to parse configuration File.
type ParseError struct {
	err error
}

// Error returns the formatted configuration error.
func (e ParseError) Error() string {
	return fmt.Sprintf("While parsing config: %s", e.err.Error())
}

// NewConfig returns a *Config with the default settings of the application.
func NewConfig(dir string) *Config {
	cfg := defaultSettings()
	absPath, err := filesystem.AbsolutePath(dir)
	if err != nil {
		absPath = dir
	}
	cfg.ConfigPath = absPath

	return cfg
}

// LoadConfigAndValidate reads configuration from the specified directory and validate it.
func LoadConfigAndValidate(dir string) (*Config, error) {
	cfg, err := LoadConfigFromFile(dir)
	if err != nil {
		return cfg, fmt.Errorf("could't read configuration: File=%s, err=%s", dir, err)
	}
	if err = cfg.Validate(); err != nil {
		return cfg, ParseError{err}
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	if len(c.Jobs) < 1 {
		return fmt.Errorf("no Jobs has been supplied")
	}

	for i := range c.Jobs {
		item := &c.Jobs[i] // we do that way to modify original object while iterating.
		srcFolder, err := filesystem.AbsolutePath(item.SourceFolder)
		if err != nil {
			return fmt.Errorf("invalid source folder. SourceFolder=%s, err=%s", item.SourceFolder, err)
		}
		item.SourceFolder = srcFolder
		if !filesystem.IsDir(item.SourceFolder) {
			return fmt.Errorf("invalid source folder. SourceFolder=%s", item.SourceFolder)
		}
	}

	return nil
}

// CompletedUploadsDBDir returns the path of the folder where completed uploads are tracked.
func (c *Config) CompletedUploadsDBDir() string {
	return path.Join(c.ConfigPath, "uploads.db")
}

// ResumableUploadsDBDir returns the path of the folder where upload URLs are tracked.
func (c *Config) ResumableUploadsDBDir() string {
	return path.Join(c.ConfigPath, "resumable_uploads.db")
}

// File return the path of the configuration File.
func (c *Config) File() string {
	return path.Join(c.ConfigPath, DefaultConfigFilename)
}

// KeyringDir returns the path of the folder where keyring will be stored.
// This is only used if 'SecretsBackendType=File'
func (c *Config) KeyringDir() string {
	return c.ConfigPath
}

// String returns a string representation of the Config object
func (c *Config) String() string {
	configTemplate := `
{
  SecretsBackendType: %s,
  APIAppCredentials: {
    ClientID:     "%s",
    ClientSecret: "%s",
  }
  jobs: [
    {
      account: %s
      sourceFolder: %s
      makeAlbums: {
        enabled: %t
        use: %s
      }
      deleteAfterUpload: %t
      includePatterns: []
	  excludePatterns: []
    }
  ]
}`
	return fmt.Sprintf(configTemplate,
		c.SecretsBackendType,
		c.APIAppCredentials.ClientID,
		c.APIAppCredentials.ClientSecret,
		c.Jobs[0].Account,
		c.Jobs[0].SourceFolder,
		c.Jobs[0].MakeAlbums.Enabled,
		c.Jobs[0].MakeAlbums.Use,
		c.Jobs[0].DeleteAfterUpload)
}

func (c *Config) WriteToFile() error {
	fh, err := os.Create(c.File())
	if err != nil {
		return err
	}
	defer fh.Close()

	_, err = fh.WriteString(c.String())
	if err != nil {
		return fmt.Errorf("failed to write configuration: File=%s, err=%v", c.File(), err)
	}

	return fh.Sync()
}

// LoadConfigFromFile reads configuration from the specified directory.
// It reads a HJSON File (given by config.File() func) and decodes it.
func LoadConfigFromFile(dir string) (*Config, error) {
	cfg := NewConfig(dir)

	file, err := ioutil.ReadFile(cfg.File())
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration: File=%s, err=%v", cfg.File(), err)
	}

	if err := unmarshalReader(bytes.NewReader(file), cfg); err != nil {
		return nil, ParseError{err}
	}

	return cfg, nil
}

// unmarshalReader unmarshal HJSON data.
func unmarshalReader(in io.Reader, c interface{}) error {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(in)

	b, err := convertHjsonToJson(buf.Bytes())
	if err != nil {
		return err
	}

	// unmarshal
	return json.Unmarshal(b, c)
}

// convertHjsonToJson converts from HJSON to JSON.
func convertHjsonToJson(in []byte) ([]byte, error) {
	var raw map[string]interface{}
	if err := hjson.Unmarshal(in, &raw); err != nil {
		return nil, err
	}

	// convert to JSON
	return json.Marshal(raw)
}

// InitConfigFile creates a config File with default settings.
func InitConfigFile(dir string) error {
	cfg := NewConfig(dir)

	if err := filesystem.EmptyOrCreateDir(cfg.ConfigPath); err != nil {
		return fmt.Errorf("failed to create config directory: path=%s, err=%v", cfg.ConfigPath, err)
	}

	return cfg.WriteToFile()
}

// Exists checks if a gphotos-uplaoder-cli configuration exists at a certain path
func Exists(path string) bool {
	cfgFile, err := filesystem.AbsolutePath(filepath.Join(path, DefaultConfigFilename))
	if err != nil {
		return false
	}

	if _, err := os.Stat(cfgFile); err == nil {
		return true
	}

	return false
}

// defaultSettings() returns a *Config with the default settings of the application.
func defaultSettings() *Config {
	return &Config{
		SecretsBackendType: "auto",
		APIAppCredentials: APIAppCredentials{
			ClientID:     "20637643488-1hvg8ev08r4tc16ca7j9oj3686lcf0el.apps.googleusercontent.com",
			ClientSecret: "0JyfLYw0kyDcJO-pGg5-rW_P",
		},
		Jobs: []FolderUploadJob{
			{
				Account:      "youremail@gmail.com",
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
