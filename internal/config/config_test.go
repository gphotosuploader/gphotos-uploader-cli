package config_test

import (
	"path/filepath"
	"testing"

	"github.com/Flaque/filet"
	"github.com/spf13/afero"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/config"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/filesystem"
)


var AppFs = afero.NewOsFs()

func TestNew(t *testing.T) {
	// Setup a temporary FS to run tests
	defer filet.CleanUp(t)
	testsFolder := filet.TmpDir(t, "")

	testCases := []struct{
		name          string
		path          string
		isErrExpected bool
	}{
		{"Should success with absolute path", testsFolder + "/foo", false},
		{"Should fail due to permission denial", "/foo", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := config.New(AppFs, tc.path)
			assertExpectedError(t, tc.isErrExpected, err)

			filename := filepath.Join(tc.path, config.DefaultConfigFilename)
			if !tc.isErrExpected && !filesystem.IsFile(filename) {
				t.Fatalf("not created: %s", filename)
			}
		})
	}
}

func assertExpectedError(t *testing.T, errExpected bool, err error) {
	if errExpected && err == nil {
		t.Fatalf("error was expected, but not produced")
	}
	if !errExpected && err != nil {
		t.Fatalf("error was not expected, err: %s", err)
	}
}

/*
func TestInitConfig(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))

	t.Run("TestInitConfigFile", func(t *testing.T) {
		err := config.New(dir)
		if err != nil {
			t.Errorf("could not create init config File: %v", err)
		}
	})
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Errorf("could not remove test config File (dir: %s): %v", dir, err)
		}
	}()
}

func TestInitAndLoadConfig(t *testing.T) {
	// init config folder
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))

	t.Run("TestInitConfigFile", func(t *testing.T) {
		err := config.New(dir)
		if err != nil {
			t.Errorf("could not create init config File: %v", err)
		}
	})
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Errorf("could not remove test config File (dir: %s): %v", dir, err)
		}
	}()

	// prepare expected configuration
	want := createTestConfiguration()

	t.Run("TestLoadConfigFile", func(t *testing.T) {
		// test load config File
		got, err := config.LoadConfigFromFile(dir)
		if err != nil {
			t.Errorf("could not load config File, got an error: %v", err)
		}

		// check that both configuration are equal
		if got.APIAppCredentials != want.APIAppCredentials {
			t.Errorf("APIAppCredentials are not equal: expected %v, got %v", want.APIAppCredentials, got.APIAppCredentials)
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
		t.Errorf("could not remove test config File (dir: %s): %v", dir, err)
	}

	_, err = config.LoadConfigFromFile(dir)
	if err == nil {
		t.Error("an error loading a non existent File was expected")
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
			t.Errorf("could not remove test config File (dir: %s): %v", dir, err)
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

func TestConfigExists(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-test.%d", time.Now().UnixNano()))
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("no error was expected at this point: err=%s", err)
	}

	t.Run("TestNonExistingConfiguration", func(t *testing.T) {
		if config.Exists(dir) {
			t.Errorf("config File exists. That was not expected: dir=%s", dir)
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
		t.Fatalf("failed to create config: File=%s, err=%v", cfgFile, err)
	}
	defer func() {
		if err := fh.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	if _, err = fh.WriteString("testTest"); err != nil {
		t.Fatalf("failed to write configuration: File=%s, err=%v", cfgFile, err)
	}

	t.Run("TestExistingConfiguration", func(t *testing.T) {
		if !config.Exists(dir) {
			t.Errorf("config File doesn't exist. That was not expected: dir=%s", dir)
		}
	})
}

func TestLoadConfigAndValidate(t *testing.T) {
	defer filet.CleanUp(t)

	want := struct {
		cfgDir    string
		srcFolder string
	}{
		cfgDir:    filet.TmpDir(t, ""),
		srcFolder: filet.TmpDir(t, ""),
	}

	t.Run("WithValidSourceFolder", func(t *testing.T) {
		// prepare a valid config File
		cfg := createTestConfiguration()
		cfg.Jobs[0].SourceFolder = want.srcFolder
		cfg.ConfigPath = want.cfgDir
		if err := cfg.writeFile(cfg.File()); err != nil {
			t.Fatal(err)
		}

		got, err := config.FromFile(want.cfgDir)
		if err != nil {
			t.Errorf("failed to load and Validate config File: err=%v", err)
		}

		if got.Jobs[0].SourceFolder != want.srcFolder {
			t.Errorf("failed: want=%s, got=%s", want.srcFolder, got.Jobs[0].SourceFolder)
		}
	})
	t.Run("WithInvalidSourceFolder", func(t *testing.T) {
		// prepare a valid config File
		cfg := createTestConfiguration()
		cfg.Jobs[0].SourceFolder = want.srcFolder
		if err := os.RemoveAll(want.srcFolder); err != nil {
			t.Fatalf("could not remove test source folder: err=%v", err)
		}
		cfg.ConfigPath = want.cfgDir
		if err := cfg.writeFile(cfg.File()); err != nil {
			t.Fatal(err)
		}

		if _, err := config.FromFile(want.cfgDir); err == nil {
			t.Errorf("failed: invalid configuration was expected")
		}
	})
	t.Run("WithNonExistentConfig", func(t *testing.T) {
		// prepare a valid config File
		cfg := createTestConfiguration()
		cfg.Jobs[0].SourceFolder = want.srcFolder
		cfg.ConfigPath = want.cfgDir
		if err := cfg.writeFile(cfg.File()); err != nil {
			t.Fatal(err)
		}
		if err := os.RemoveAll(want.cfgDir); err != nil {
			t.Fatalf("could not remove config folder: err=%v", err)
		}

		if _, err := config.FromFile(want.cfgDir); err == nil {
			t.Errorf("failed: invalid configuration was expected")
		}
	})
}

func createTestConfiguration() *config.AppConfig {
	fc := &config.Config{
		SecretsBackendType: "auto",
		APIAppCredentials: config.APIAppCredentials{
			ClientID:     "20637643488-1hvg8ev08r4tc16ca7j9oj3686lcf0el.apps.googleusercontent.com",
			ClientSecret: "0JyfLYw0kyDcJO-pGg5-rW_P",
		},
		Account: "youremail@gmail.com",
		Jobs: []config.FolderUploadJob{
			{
				SourceFolder: "~/folder/to/upload",
				MakeAlbums: config.MakeAlbums{
					Enabled: true,
					Use:     "folderName",
				},
				DeleteAfterUpload: false,
			},
		},
	}
	return &config.AppConfig{
		Config: fc,
	}
}
*/