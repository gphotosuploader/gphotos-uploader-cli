package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/client9/xson/hjson"
	"github.com/nmrshll/go-cp"

	"github.com/gphotosuploader/gphotos-uploader-cli/utils/filesystem"
)

// defaultSettings() returns a *Config with the default settings of the application.
func defaultSettings() *Config {
	var c Config
	c.SecretsBackendType = "auto"
	c.APIAppCredentials = &APIAppCredentials{
		ClientID:     "20637643488-1hvg8ev08r4tc16ca7j9oj3686lcf0el.apps.googleusercontent.com",
		ClientSecret: "0JyfLYw0kyDcJO-pGg5-rW_P",
	}
	c.Jobs = make([]FolderUploadJob, 0)
	job := FolderUploadJob{
		Account:      "youremail@gmail.com",
		SourceFolder: "~/folder/to/upload",
		MakeAlbums: MakeAlbums{
			Enabled: true,
			Use:     "folderNames",
		},
		DeleteAfterUpload: true,
		UploadVideos:      true,
	}
	c.Jobs = append(c.Jobs, job)
	return &c
}

// NewConfig returns a *Config with the default settings of the application.
func NewConfig(dir string) *Config {
	cfg := defaultSettings()
	cfg.ConfigPath = filesystem.AbsolutePath(dir)

	return cfg
}

func (c *Config) Validate() error {
	if len(c.Jobs) < 1 {
		return fmt.Errorf("no Jobs has been supplied")
	}

	for _, item := range c.Jobs {
		path, err := cp.AbsolutePath(item.SourceFolder)
		if err != nil {
			return fmt.Errorf("invalid source folder. SourceFolder=%s, err=%s", item.SourceFolder, err)
		}
		item.SourceFolder = path
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

// ConfigFile return the path of the configuration file.
func (c *Config) ConfigFile() string {
	return path.Join(c.ConfigPath, "config.hjson")
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
      uploadVideos: %t
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
		c.Jobs[0].DeleteAfterUpload,
		c.Jobs[0].UploadVideos)
}

// LoadConfig reads configuration from the specified directory.
// It reads a HJSON file (given by config.ConfigFile() func) and decodes it.
func LoadConfig(dir string) (*Config, error) {
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

// InitConfig creates an example config file if it doesn't already exist.
// If 'force' is set then we are going to remove config dir before creating it.
func InitConfig(dir string, force bool) error {
	cfg := NewConfig(dir)

	// if force, we should remove everything to start from the scratch.
	if force {
		err := os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(cfg.ConfigPath); !os.IsNotExist(err) {
		// directory already exist and forced was not set
		return fmt.Errorf("config directory already exists, use '--force' to overwrite: path=%s", cfg.ConfigPath)
	} else {
		err := os.MkdirAll(cfg.ConfigPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create config directory: path=%s, err=%v", cfg.ConfigPath, err)
		}
	}

	fh, err := os.Open(cfg.ConfigFile())
	if err != nil {
		fh, err = os.Create(cfg.ConfigFile())
		if err != nil {
			return fmt.Errorf("failed to create config: file=%s, err=%v", cfg.ConfigFile(), err)
		}
	}
	defer func() {
		if err := fh.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	_, err = fh.WriteString(cfg.String())
	if err != nil {
		return fmt.Errorf("failed to write configuration: file=%s, err=%v", cfg.ConfigFile(), err)
	}

	return fh.Sync()
}
