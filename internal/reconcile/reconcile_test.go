package reconcile_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/config"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/filetracker"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/reconcile"
)

func TestBuildDetectsPendingAndChangedBytes(t *testing.T) {
	dir := t.TempDir()
	write := func(name, content string) string {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	accepted := write("accepted.jpg", "original bytes")
	pending := write("pending.mov", "pending")
	repo, err := filetracker.NewLevelDBRepository(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	tracker := filetracker.New(repo)
	defer tracker.Close()
	hash, err := (filetracker.SHA256Hasher{}).Hash(accepted)
	if err != nil {
		t.Fatal(err)
	}
	if err := tracker.RecordUpload(accepted, filetracker.UploadReceipt{MediaItemID: "google-id", SHA256: hash, Size: int64(len("original bytes")), UploadedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{Jobs: []config.FolderUploadJob{{SourceFolder: dir, IncludePatterns: []string{"_ALL_FILES_"}}}}
	report, err := reconcile.Build(cfg, tracker)
	if err != nil {
		t.Fatal(err)
	}
	if report.NASItems != 2 || report.VerifiedItems != 1 || report.PendingItems != 1 || report.GoogleAcceptedItems != 1 {
		t.Fatalf("unexpected report: %+v", report)
	}
	if report.PendingBytes != int64(len("pending")) {
		t.Fatalf("pending bytes = %d", report.PendingBytes)
	}
	if err := os.WriteFile(accepted, []byte("changed bytes!"), 0600); err != nil {
		t.Fatal(err)
	}
	report, err = reconcile.Build(cfg, tracker)
	if err != nil {
		t.Fatal(err)
	}
	if report.MismatchedItems != 1 || report.VerifiedItems != 0 {
		t.Fatalf("changed content was not detected: %+v", report)
	}
	_ = pending
}
