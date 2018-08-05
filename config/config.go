package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	cp "github.com/nmrshll/go-cp"
	"github.com/nmrshll/gphotos-uploader-cli/fileshandling"
	"github.com/palantir/stacktrace"

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
	CONFIGPATH          = fmt.Sprintf("%s/config.hjson", CONFIGFOLDER)
	API_APP_CREDENTIALS = APIAppCredentials{
		ClientID:     "20637643488-1hvg8ev08r4tc16ca7j9oj3686lcf0el.apps.googleusercontent.com",
		ClientSecret: "0JyfLYw0kyDcJO-pGg5-rW_P",
	}

	// vars
	Cfg *Config
)

type FolderUploadJob struct {
	// Credentials  *auth.CookieCredentials
	Account      string
	SourceFolder string
	MakeAlbums   struct {
		Enabled bool
		Use     string
	}
	DeleteAfterUpload bool
}
type Config struct {
	Accounts map[string]struct {
		Username string
		Password string
	}
	Jobs []*FolderUploadJob
}

func Load() *Config {
	Cfg = loadConfigFile()
	return Cfg
}

// func processConfig(jobsConfig *Config) {
// 	for _, job := range jobsConfig.Jobs {
// 		authFilePathRaw := fmt.Sprintf("%s/%s/auth.json", CONFIGFOLDER, slug.Make(job.Account))
// 		authFilePathAbsolute, err := cp.AbsolutePath(authFilePathRaw)
// 		_ = authFilePathAbsolute
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(fmt.Sprintf("loading auth file at %s", authFilePathRaw))
// 		// job.Credentials = Authenticate(authFilePathAbsolute)
// 	}
// }

func loadConfigFile() *Config {
	configPathAbsolute, err := cp.AbsolutePath(CONFIGPATH)
	if err != nil {
		log.Fatal(err)
	}

	if fileshandling.IsFile(configPathAbsolute) {
		fmt.Println("[INFO] Config file found. Loading...")
	} else {
		err := cp.CopyFile("config.hjson.example", configPathAbsolute)
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
