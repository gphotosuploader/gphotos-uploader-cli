package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/davecgh/go-spew/spew"
	cp "github.com/nmrshll/go-cp"
	gphotos "github.com/nmrshll/google-photos-api-client-go/lib-gphotos"
	"github.com/nmrshll/gphotos-uploader-cli/fileshandling"
	"github.com/palantir/stacktrace"
	"golang.org/x/oauth2"

	"github.com/client9/xson/hjson"
)

const (
	CONFIGFOLDER = "~/.config/gphotos-uploader-go-api"
)

type APIAppCredentials struct {
	ClientID     string
	ClientSecret string
}

var (
	// consts
	CONFIGPATH                  = fmt.Sprintf("%s/config.hjson", CONFIGFOLDER)
	DEFAULT_API_APP_CREDENTIALS = APIAppCredentials{
		ClientID:     "20637643488-1hvg8ev08r4tc16ca7j9oj3686lcf0el.apps.googleusercontent.com",
		ClientSecret: "0JyfLYw0kyDcJO-pGg5-rW_P",
	}

	// vars
	Cfg *Config
)

type Config struct {
	APIAppCredentials *APIAppCredentials
	Jobs              []FolderUploadJob
}

func (c *Config) Process() {
	if c.APIAppCredentials == nil {
		c.APIAppCredentials = &DEFAULT_API_APP_CREDENTIALS
	}
}

func OAuthConfig() *oauth2.Config {
	if Cfg.APIAppCredentials == nil {
		log.Fatal(stacktrace.NewError("APIAppCredentials can't be nil"))
	}
	return gphotos.NewOAuthConfig(gphotos.APIAppCredentials(*Cfg.APIAppCredentials))
}

type FolderUploadJob struct {
	Account      string
	SourceFolder string
	MakeAlbums   struct {
		Enabled bool
		Use     string
	}
	DeleteAfterUpload bool
}

func Load() *Config {
	Cfg = loadConfigFile()
	Cfg.Process()
	spew.Dump(Cfg)
	return Cfg
}

func loadConfigFile() *Config {
	configPathAbsolute, err := cp.AbsolutePath(CONFIGPATH)
	if err != nil {
		log.Fatal(err)
	}

	if fileshandling.IsFile(configPathAbsolute) {
		fmt.Println("[INFO] Config file found. Loading...")
	} else {
		err := cp.CopyFile("./config/config.example.hjson", configPathAbsolute)
		if err != nil {
			log.Fatal(err)
		}
	}

	configBytes, err := ioutil.ReadFile(configPathAbsolute)
	if err != nil {
		log.Fatal(err)
	}
	var config = &Config{}
	jsonConfig := hjson.ToJSON(configBytes)
	if err := json.Unmarshal(jsonConfig, config); err != nil {
		log.Fatal(stacktrace.Propagate(err, "unmarshal jsonConfig failed"))
	}
	return config
}
