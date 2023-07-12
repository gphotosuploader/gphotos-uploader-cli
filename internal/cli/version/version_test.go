package version_test

import (
	"bytes"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/version"
	"io"
	"testing"
)

func TestNewCommand(t *testing.T) {
	c := version.NewCommand()
	b := bytes.NewBufferString("")
	c.SetOut(b)
	if err := c.Execute(); err != nil {
		t.Fatal(err)
	}
	got, err := io.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	want := "gphotos-cli Version: 0.0.0-git\n"
	if want != string(got) {
		t.Fatalf("want: %s, got: %s", want, string(got))
	}
}
