package config_test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/gphotosuploader/gphotos-uploader-cli/config"
)

func TestInitConfig(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))

	t.Run("TestInitConfigFile", func(t *testing.T) {
		err := config.InitConfig(dir, true)
		if err != nil {
			t.Errorf("could not create init config file: %v", err)
		}
	})

	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Errorf("could not remove test config file (dir: %s): %v", dir, err)
		}
	}()

	t.Run("TestInitConfigFileWithExistentFile", func(t *testing.T) {
		err := config.InitConfig(dir, false)
		if err == nil {
			t.Error("an error creating an existent file was expected")
		}
	})
}

func TestInitAndLoadConfig(t *testing.T) {
	// init config folder
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))

	t.Run("TestInitConfigFile", func(t *testing.T) {
		err := config.InitConfig(dir, true)
		if err != nil {
			t.Errorf("could not create init config file: %v", err)
		}
	})

	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Errorf("could not remove test config file (dir: %s): %v", dir, err)
		}
	}()

	// prepare expected configuration
	expected := createTestConfiguration()

	t.Run("TestLoadConfigFile", func(t *testing.T) {
		// test load config file
		got, err := config.LoadConfig(dir)
		if err != nil {
			t.Errorf("could not load config file, got an error: %v", err)
		}

		// check that both configuration are equal
		if *got.APIAppCredentials != *expected.APIAppCredentials {
			t.Errorf("APIAppCredentials are not equal: expected %v, got %v", *expected.APIAppCredentials, *got.APIAppCredentials)
		}

		if len(got.Jobs) != len(expected.Jobs) {
			t.Errorf("Jobs are not equal: expected %d jobs, got %d jobs", len(expected.Jobs), len(got.Jobs))
		}
	})
}

func TestLoadConfigWithNonExistentFile(t *testing.T) {
	// init config folder
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))
	err := os.RemoveAll(dir)
	if err != nil {
		t.Errorf("could not remove test config file (dir: %s): %v", dir, err)
	}

	_, err = config.LoadConfig(dir)
	if err == nil {
		t.Error("an error loading a non existent file was expected")
	}
}

func TestConfig_CompletedUploadsDBDir(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))
	c := config.NewConfig(dir)

	expected := path.Join(dir, "uploads.db")
	got := c.CompletedUploadsDBDir()

	if got != expected {
		t.Errorf("Testing get completed uploads DB dir: expected: %s, got %s", expected, got)
	}

}

func TestConfig_ResumableUploadsDBDir(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))
	c := config.NewConfig(dir)

	expected := path.Join(dir, "resumable_uploads.db")
	got := c.ResumableUploadsDBDir()

	if got != expected {
		t.Errorf("Testing get resumable uploads DB dir: expected: %s, got %s", expected, got)
	}

}

func TestConfig_KeyringDir(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))
	c := config.NewConfig(dir)

	expected := dir
	got := c.KeyringDir()

	if got != expected {
		t.Errorf("Testing get keyring dir: expected: %s, got %s", expected, got)
	}
}

func createTestConfiguration() *config.Config {
	c := &config.Config{}
	c.SecretsBackendType = "auto"
	c.APIAppCredentials = &config.APIAppCredentials{
		ClientID:     "20637643488-1hvg8ev08r4tc16ca7j9oj3686lcf0el.apps.googleusercontent.com",
		ClientSecret: "0JyfLYw0kyDcJO-pGg5-rW_P",
	}
	c.Jobs = make([]config.FolderUploadJob, 0)
	job := config.FolderUploadJob{
		Account:      "youremail@gmail.com",
		SourceFolder: "~/folder/to/upload",
		MakeAlbums: config.MakeAlbums{
			Enabled: true,
			Use:     "folderNames",
		},
		DeleteAfterUpload: true,
		UploadVideos:      true,
	}
	c.Jobs = append(c.Jobs, job)
	return c
}
