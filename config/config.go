package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"gitlab.com/nmrshll/gphotos-uploader-go-cookies/filesHandling"
	"gitlab.com/nmrshll/gphotos-uploader-go-cookies/upload"

	"github.com/gosimple/slug"
	cp "github.com/nmrshll/go-cp"
	"github.com/palantir/stacktrace"

	"github.com/client9/xson/hjson"
)

const (
	CONFIGFOLDER = "~/.config/gphotos-uploader-go-cookies"
)

var (
	CONFIGPATH = fmt.Sprintf("%s/config.hjson", CONFIGFOLDER)
)

type JobsConfig struct {
	Accounts map[string]struct {
		Username string
		Password string
	}
	Jobs []*upload.FolderUploadJob
}

func Load() []*upload.FolderUploadJob {
	jobConfig := loadConfigFile()
	for _, job := range jobConfig.Jobs {
		authFilePath := fmt.Sprintf("%s/%s/auth.json", CONFIGFOLDER, slug.Make(job.Account))
		authFilePathAbsolute, err := cp.AbsolutePath(authFilePath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(fmt.Sprintf("loading auth file at %s", authFilePath))
		job.Credentials = Authenticate(authFilePathAbsolute)
	}
	return jobConfig.Jobs
}

func loadConfigFile() *JobsConfig {
	configPathAbsolute, err := cp.AbsolutePath(CONFIGPATH)
	if err != nil {
		log.Fatal(err)
	}

	if filesHandling.IsFile(configPathAbsolute) {
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
	var config = &JobsConfig{}
	jsonConfig := hjson.ToJSON(configBytes)
	if err := json.Unmarshal(jsonConfig, config); err != nil {
		log.Fatal(stacktrace.Propagate(err, "unmarshal jsonConfig failed"))
	}
	return config
}
