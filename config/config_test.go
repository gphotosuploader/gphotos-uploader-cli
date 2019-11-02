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
		err := config.InitConfigFile(dir)
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
}

func TestInitAndLoadConfig(t *testing.T) {
	// init config folder
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))

	t.Run("TestInitConfigFile", func(t *testing.T) {
		err := config.InitConfigFile(dir)
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
	want := createTestConfiguration()

	t.Run("TestLoadConfigFile", func(t *testing.T) {
		// test load config file
		got, err := config.LoadConfig(dir)
		if err != nil {
			t.Errorf("could not load config file, got an error: %v", err)
		}

		// check that both configuration are equal
		if *got.APIAppCredentials != *want.APIAppCredentials {
			t.Errorf("APIAppCredentials are not equal: expected %v, got %v", *want.APIAppCredentials, *got.APIAppCredentials)
		}

		if len(got.Jobs) != len(want.Jobs) {
			t.Errorf("Jobs are not equal: expected %d jobs, got %d jobs", len(want.Jobs), len(got.Jobs))
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

	want := path.Join(dir, "uploads.db")
	got := c.CompletedUploadsDBDir()

	if got != want {
		t.Errorf("Testing get completed uploads DB dir: expected: %s, got %s", want, got)
	}

}

func TestConfig_ResumableUploadsDBDir(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))
	c := config.NewConfig(dir)

	want := path.Join(dir, "resumable_uploads.db")
	got := c.ResumableUploadsDBDir()

	if got != want {
		t.Errorf("Testing get resumable uploads DB dir: expected: %s, got %s", want, got)
	}

}

func TestConfig_KeyringDir(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))
	c := config.NewConfig(dir)

	want := dir
	got := c.KeyringDir()

	if got != want {
		t.Errorf("Testing get keyring dir: expected: %s, got %s", want, got)
	}
}

func TestConfig_Validate(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-test.%d", time.Now().UnixNano()))
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		t.Errorf("no error was expected at this point: err=%s", err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Errorf("could not remove test config file (dir: %s): %v", dir, err)
		}
	}()

	t.Run("TestValidateConfigWithValidSettings", func(t *testing.T) {
		c := createTestConfiguration()
		c.Jobs[0].SourceFolder = dir
		got := c.Validate()
		if got != nil {
			t.Errorf("config is not validated. That was not expected: err=%v", got)
		}
	})

	t.Run("TestValidateConfigWithInvalidSourceFolder", func(t *testing.T) {
		c := createTestConfiguration()
		got := c.Validate()
		if got == nil {
			t.Errorf("config is validated. That was not expected: err=%v", got)
		}
	})

	t.Run("TestValidateConfigWithoutJobs", func(t *testing.T) {
		c := createTestConfiguration()
		c.Jobs = make([]config.FolderUploadJob, 0)
		got := c.Validate()
		if got == nil {
			t.Errorf("config is validated. That was not expected: err=%v", got)
		}
	})

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

func TestConfigExists(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-test.%d", time.Now().UnixNano()))
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("no error was expected at this point: err=%s", err)
	}

	t.Run("TestNonExistingConfiguration", func(t *testing.T) {
		if config.ConfigExists(dir) {
			t.Errorf("config file exists. That was not expected: dir=%s", dir)
		}
	})

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		t.Fatalf("no error was expected at this point: err=%s", err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatalf("could not remove test config folder: path=%s err=%s", dir, err)
		}
	}()

	cfgFile := filepath.Join(dir, "config.hjson")
	fh, err := os.Create(cfgFile)
	if err != nil {
		t.Fatalf("failed to create config: file=%s, err=%v", cfgFile, err)
	}
	defer func() {
		if err := fh.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	if _, err = fh.WriteString("testTest"); err != nil {
		t.Fatalf("failed to write configuration: file=%s, err=%v", cfgFile, err)
	}

	t.Run("TestExistingConfiguration", func(t *testing.T) {
		if !config.ConfigExists(dir) {
			t.Errorf("config file doesn't exist. That was not expected: dir=%s", dir)
		}
	})
}
