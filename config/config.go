package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/nmrshll/go-cp"
	"github.com/nmrshll/google-photos-api-client-go/lib-gphotos"
	"github.com/palantir/stacktrace"
	"golang.org/x/oauth2"

	"github.com/client9/xson/hjson"
)

const exampleConfig = `
{
  APIAppCredentials: {
    ClientID:     "20637643488-1hvg8ev08r4tc16ca7j9oj3686lcf0el.apps.googleusercontent.com",
    ClientSecret: "0JyfLYw0kyDcJO-pGg5-rW_P",
  }
  jobs: [
    {
      account: youremail@gmail.com
      sourceFolder: ~/folder/to/upload
      makeAlbums: {
        enabled: true
        use: folderNames
      }
      deleteAfterUpload: true
    }
  ]
}
`

type APIAppCredentials struct {
	ClientID     string
	ClientSecret string
}

type FolderUploadJob struct {
	Account           string
	SourceFolder      string
	MakeAlbums        MakeAlbums
	DeleteAfterUpload bool
}

type MakeAlbums struct {
	Enabled bool
	Use     string
}

type Config struct {
	APIAppCredentials *APIAppCredentials
	Jobs              []FolderUploadJob
}

var (
	Cfg *Config
)

func OAuthConfig() *oauth2.Config {
	if Cfg.APIAppCredentials == nil {
		log.Fatal(stacktrace.NewError("APIAppCredentials can't be nil"))
	}
	return gphotos.NewOAuthConfig(gphotos.APIAppCredentials(*Cfg.APIAppCredentials))
}

// LoadConfigFile reads HJSON file (absolute path) and decodes its configuration
func LoadConfigFile(p string) (*Config, error) {
	path, err := cp.AbsolutePath(p)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to get absolute path: %s", p)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to read configuration from file: %s", path)
	}

	var config = &Config{}

	if err := hjson.Unmarshal(data, config); err != nil {
		return nil, stacktrace.Propagate(err, "Failed to decode configuration data")
	}

	return config, nil
}

// InitConfigFile creates an example config file if it doesn't already exist
func InitConfigFile(p string) error {
	path, err := cp.AbsolutePath(p)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to get absolute path: %s", p)
	}

	dirname := filepath.Dir(path)
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		err := os.MkdirAll(dirname, 0755)
		if err != nil {
			return stacktrace.Propagate(err, "Failed to create config parent directory: %s", dirname)
		}
	}

	fh, err := os.Open(path)
	if err != nil {
		fh, err = os.Create(path)
		if err != nil {
			return stacktrace.Propagate(err, "Failed to create config file: %s", path)
		}
	}
	defer fh.Close()

	_, err = fh.WriteString(exampleConfig)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to write in config file: %s", path)
	}

	return fh.Sync()
}
