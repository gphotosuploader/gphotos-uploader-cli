package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/nmrshll/go-cp"
	"github.com/nmrshll/google-photos-api-client-go/lib-gphotos"
	"golang.org/x/oauth2"

	"github.com/client9/xson/hjson"
)

type APIAppCredentials struct {
	ClientID     string
	ClientSecret string
}

type FolderUploadJob struct {
	Account           string
	SourceFolder      string
	MakeAlbums        MakeAlbums
	DeleteAfterUpload bool
	UploadVideos      bool
}

type MakeAlbums struct {
	Enabled bool
	Use     string
}

type Config struct {
	APIAppCredentials *APIAppCredentials
	Jobs              []FolderUploadJob
}

// defaultConfig returns an example Config object
func defaultConfig() *Config {
	c := &Config{}
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
		UploadVideos:      false,
	}
	c.Jobs = append(c.Jobs, job)
	return c
}

// String returns a string representation of the Config object
func (c Config) String() string {
	configTemplate := `
{
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
    }
  ]
}`
	return fmt.Sprintf(configTemplate,
		c.APIAppCredentials.ClientID,
		c.APIAppCredentials.ClientSecret,
		c.Jobs[0].Account,
		c.Jobs[0].SourceFolder,
		c.Jobs[0].MakeAlbums.Enabled,
		c.Jobs[0].MakeAlbums.Use,
		c.Jobs[0].DeleteAfterUpload,
		c.Jobs[0].UploadVideos)
}

var (
	Cfg *Config
)

func OAuthConfig() *oauth2.Config {
	if Cfg.APIAppCredentials == nil {
		log.Fatal(fmt.Errorf("APIAppCredentials can't be nil"))
	}
	return gphotos.NewOAuthConfig(gphotos.APIAppCredentials(*Cfg.APIAppCredentials))
}

// GetUploadsDBPath returns the absolute path of uploads DB file
func GetUploadsDBPath() string {
	const UploadDBPath = "~/.config/gphotos-uploader-cli/uploads.db"

	dbPathAbsolute, err := cp.AbsolutePath(UploadDBPath)
	if err != nil {
		log.Fatal(err) // TODO: should return an error instead a log.Fatal
	}
	return dbPathAbsolute
}

// LoadConfigFile reads HJSON file (absolute path) and decodes its configuration
func LoadConfigFile(p string) (*Config, error) {
	path, err := cp.AbsolutePath(p)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %s", p)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration from file %s, %v", path, err)
	}

	var config = &Config{}

	if err := hjson.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to decode configuration data: %v", err)
	}

	return config, nil
}

// InitConfigFile creates an example config file if it doesn't already exist
func InitConfigFile(p string) error {
	path, err := cp.AbsolutePath(p)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %s", p)
	}

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
	defer fh.Close()

	_, err = fh.WriteString(defaultConfig().String())
	if err != nil {
		return fmt.Errorf("failed to write in config file %s: %v", path, err)
	}

	return fh.Sync()
}
