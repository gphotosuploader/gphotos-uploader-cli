package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/client9/xson/hjson"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/utils/filesystem"
)

const (
	// DefaultConfigFilename is the default config file name to use
	DefaultConfigFilename = "config.hjson"
)

// defaultSettings() returns a *Config with the default settings of the application.
func defaultSettings() *Config {
	var c Config
	c.SecretsBackendType = "auto"
	c.APIAppCredentials = APIAppCredentials{
		ClientID:     "20637643488-1hvg8ev08r4tc16ca7j9oj3686lcf0el.apps.googleusercontent.com",
		ClientSecret: "0JyfLYw0kyDcJO-pGg5-rW_P",
	}
	c.Jobs = make([]FolderUploadJob, 0)
	job := FolderUploadJob{
		Account:      "youremail@gmail.com",
		SourceFolder: "~/folder/to/upload",
		MakeAlbums: MakeAlbums{
			Enabled: true,
			Use:     "folderName",
		},
		DeleteAfterUpload: false,
	}
	c.Jobs = append(c.Jobs, job)
	return &c
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

func (c *Config) Validate() error {
	if len(c.Jobs) < 1 {
		return fmt.Errorf("no Jobs has been supplied")
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

// ConfigFile return the path of the configuration file.
func (c *Config) ConfigFile() string {
	return path.Join(c.ConfigPath, DefaultConfigFilename)
}

// KeyringDir returns the path of the folder where keyring will be stored.
// This is only used if 'SecretsBackendType=file'
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
	fh, err := os.Create(c.ConfigFile())
	if err != nil {
		return err
	}
	defer fh.Close()

	_, err = fh.WriteString(c.String())
	if err != nil {
		return fmt.Errorf("failed to write configuration: file=%s, err=%v", c.ConfigFile(), err)
	}

	return fh.Sync()
}

// LoadConfigFromFile reads configuration from the specified directory.
// It reads a HJSON file (given by config.ConfigFile() func) and decodes it.
func LoadConfigFromFile(dir string) (*Config, error) {
	cfg := NewConfig(dir)

	data, err := ioutil.ReadFile(cfg.ConfigFile())
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration: file=%s, err=%v", cfg.ConfigFile(), err)
	}

	if err := hjson.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode configuration data: err=%v", err)
	}

	return cfg, nil
}

// LoadConfigAndValidate reads configuration from the specified directory and validate it.
func LoadConfigAndValidate(dir string) (*Config, error) {
	cfg, err := LoadConfigFromFile(dir)
	if err != nil {
		return cfg, fmt.Errorf("could't read configuration: file=%s, err=%s", dir, err)
	}
	if err = cfg.Validate(); err != nil {
		return cfg, fmt.Errorf("invalid configuration: file=%s, err=%s", cfg.ConfigFile(), err)
	}
	return cfg, nil
}

// InitConfigFile creates a config file with default settings.
func InitConfigFile(dir string) error {
	cfg := NewConfig(dir)

	if err := filesystem.EmptyOrCreateDir(cfg.ConfigPath); err != nil {
		return fmt.Errorf("failed to create config directory: path=%s, err=%v", cfg.ConfigPath, err)
	}

	return cfg.WriteToFile()
}

// ConfigExists checks if a gphotos-uplaoder-cli configuration exists at a certain path
func ConfigExists(path string) bool {
	cfgFile, err := filesystem.AbsolutePath(filepath.Join(path, DefaultConfigFilename))
	if err != nil {
		return false
	}

	if _, err := os.Stat(cfgFile); err == nil {
		return true
	}

	return false
}
