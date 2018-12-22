package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInitAndLoadConfig(t *testing.T) {
	// init config file
	path := filepath.Join(os.TempDir(), fmt.Sprintf("file.%d", time.Now().UnixNano()))

	t.Run("TestInitConfigFile", func(t *testing.T) {
		err := InitConfigFile(path)
		if err != nil {
			t.Errorf("could not create init config file: %v", err)
		}
	})

	defer func() {
		err := os.Remove(path)
		if err != nil {
			t.Errorf("could not remove test config file (path: %s): %v", path, err)
		}
	}()

	// prepare expected configuration
	expected := newExampleConfig()

	t.Run("TestLoadConfigFile", func(t *testing.T) {
		// test load config file
		got, err := LoadConfigFile(path)
		if err != nil {
			t.Errorf("could not load config file, got an error: %v", err)
		}

		// check that both configuration are equal
		if *got.APIAppCredentials != *expected.APIAppCredentials {
			t.Errorf("APIAppCredentials are not equal: expected %v, got %v", *expected.APIAppCredentials, *got.APIAppCredentials)
		}

		if got.Jobs[0] != expected.Jobs[0] {
			t.Errorf("Jobs are not equal: expected %v, got %v", expected.Jobs[0], got.Jobs[0])
		}
	})
}
