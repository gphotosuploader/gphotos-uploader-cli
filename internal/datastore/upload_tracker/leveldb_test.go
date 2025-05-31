package upload_tracker_test

import (
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/upload_tracker"
	"os"
	"testing"
)

func RemoveDB(path string) {
	_ = os.RemoveAll(path)
}

func TestNewStore(t *testing.T) {
	t.Run("Should success when folder is writable", func(t *testing.T) {
		name, err := os.MkdirTemp(os.TempDir(), "upload_tracker")
		if err != nil {
			t.Fatalf("error was not expected at this time: %v", err)
		}
		defer RemoveDB(name)

		store, err := upload_tracker.NewStore(name)
		if err != nil {
			t.Errorf("error was not expected: %v", err)
		}
		store.Close()
	})

	t.Run("Should fail when folder is not writable", func(t *testing.T) {
		name := "/non-existent"

		store, err := upload_tracker.NewStore(name)
		if err == nil {
			store.Close()
			t.Errorf("error was expected but not produced")
		}
	})
}

func TestLevelDBStore_GetSet(t *testing.T) {
	t.Run("Should get the value when the key is present", func(t *testing.T) {
		name, err := os.MkdirTemp(os.TempDir(), "upload_tracker")
		if err != nil {
			t.Fatalf("error was not expected at this time: %v", err)
		}
		defer RemoveDB(name)

		store, err := upload_tracker.NewStore(name)
		if err != nil {
			t.Fatalf("error was not expected at this time: %v", err)
		}
		defer store.Close()

		store.Set("fooKey", "fooValue")

		got, found := store.Get("fooKey")

		if !found || got != "fooValue" {
			t.Errorf("want: %s, got: %s", "fooValue", got)
		}
	})

	t.Run("Should return false if the key is not present", func(t *testing.T) {
		name, err := os.MkdirTemp(os.TempDir(), "upload_tracker")
		if err != nil {
			t.Fatalf("error was not expected at this time: %v", err)
		}
		defer RemoveDB(name)

		store, err := upload_tracker.NewStore(name)
		if err != nil {
			t.Fatalf("error was not expected at this time: %v", err)
		}
		defer store.Close()

		got, found := store.Get("non-existent")

		if found {
			t.Errorf("key was not expected, got: %s", got)
		}
	})
}

func TestLevelDBStore_Delete(t *testing.T) {
	t.Run("Should delete a key", func(t *testing.T) {
		name, err := os.MkdirTemp(os.TempDir(), "upload_tracker")
		if err != nil {
			t.Fatalf("error was not expected at this time: %v", err)
		}
		defer RemoveDB(name)

		store, err := upload_tracker.NewStore(name)
		if err != nil {
			t.Fatalf("error was not expected at this time: %v", err)
		}
		defer store.Close()

		store.Set("fooKey", "fooValue")

		store.Delete("fooKey")

		got, found := store.Get("fooKey")

		if found {
			t.Errorf("key was not expected, got: %s", got)
		}
	})
}
