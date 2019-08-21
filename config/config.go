package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/oauth2"

	"github.com/client9/xson/hjson"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"
	"github.com/gphotosuploader/gphotos-uploader-cli/utils/filesystem"
)

// APIAppCredentials represents Google Photos API credentials for OAuth
type APIAppCredentials struct {
	ClientID     string
	ClientSecret string
}

// FolderUploadJob represents configuration for a folder to be uploaded
type FolderUploadJob struct {
	Account           string
	SourceFolder      string
	MakeAlbums        MakeAlbums
	DeleteAfterUpload bool
	UploadVideos      bool
	IncludePatterns   []string
	ExcludePatterns   []string
}

// MakeAlbums represents configuration about how to create Albums in Google Photos
type MakeAlbums struct {
	Enabled bool
	Use     string
}

// Config represents this application configuration
type Config struct {
	ConfigPath         string
	TrackingDBPath     string
	Verbose            bool
	SecretsBackendType string
	APIAppCredentials  *APIAppCredentials
	Jobs               []FolderUploadJob
}

// defaultConfig returns an example Config object
func defaultConfig() Config {
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
	return c
}

// String returns a string representation of the Config object
func (c Config) String() string {
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

// OAuthConfig creates and returns a new oauth Config based on API app credentials found in the uploader's config file
func OAuthConfig(uploaderConfigAPICredentials *APIAppCredentials) *oauth2.Config {
	if uploaderConfigAPICredentials == nil {
		log.Fatalf("APIAppCredentials can't be nil")
	}
	return gphotos.NewOAuthConfig(gphotos.APIAppCredentials(*uploaderConfigAPICredentials))
}

// LoadConfigFile reads HJSON file (absolute path) and decodes its configuration
func LoadConfigFile(p string) (*Config, error) {
	config := defaultConfig()
	config.ConfigPath = filesystem.AbsolutePath(p)
	config.TrackingDBPath = path.Join(config.ConfigPath, "uploads.db")

	path := path.Join(config.ConfigPath, "config.hjson")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration from file %s, %v", path, err)
	}

	if err := hjson.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to decode configuration data: %v", err)
	}

	return &config, nil
}

// InitConfigFile creates an example config file if it doesn't already exist
func InitConfigFile(p string) error {
	path := filesystem.AbsolutePath(p)
	dirname := filepath.Dir(path)
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		err := os.MkdirAll(dirname, 0755)
		if err != nil {
			return fmt.Errorf("failed to create config parent directory %s: %v", dirname, err)
		}
	}

	fh, err := os.Open(path)
	if err != nil {
		fh, err = os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create config file %s: %v", path, err)
		}
	}
	defer func() {
		if err := fh.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	_, err = fh.WriteString(defaultConfig().String())
	if err != nil {
		return fmt.Errorf("failed to write in config file %s: %v", path, err)
	}

	return fh.Sync()
}
