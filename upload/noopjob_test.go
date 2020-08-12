package upload_test

import (
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/upload"
)

func TestNoOpJob_ID(t *testing.T) {
	j := upload.NoOpJob{}

	want := "noop"
	got := j.ID()

	if got != want {
		t.Errorf("want: %s, got: %s", want, got)
	}
}

func TestNoOpJob_Process(t *testing.T) {
	j := upload.NoOpJob{}

	err := j.Process()
	if err != nil {
		t.Errorf("no error was expected at this point: err=%s", err)
	}
}