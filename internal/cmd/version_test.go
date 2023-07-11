package cmd_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cmd"
)

func TestNewVersionCmd(t *testing.T) {
	c := cmd.NewVersionCmd()
	b := bytes.NewBufferString("")
	c.SetOut(b)
	if err := c.Execute(); err != nil {
		t.Fatal(err)
	}
	got, err := io.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	want := "gphotos-cli v0.0.0\n"
	if want != string(got) {
		t.Fatalf("want: %s, got: %s", want, string(got))
	}
}
