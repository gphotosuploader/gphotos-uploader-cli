package upload

import (
	"io"
	"path"

	"github.com/spf13/afero"
)

var (
	// fs represents the filesystem to use. By default it uses functions based on `os` package.
	// In testing, it uses a memory file system.
	appFS = afero.NewOsFs()
)

// FileItem represents a local file.
type FileItem struct {
	Path      string
	AlbumName string
}

func NewFileItem(path string) FileItem {
	return FileItem{
		Path: path,
	}
}

// Open returns a stream.
// Caller should close it finally.
func (m FileItem) Open() (io.ReadSeeker, int64, error) {
	f, err := appFS.Stat(m.Path)
	if err != nil {
		return nil, 0, err
	}
	r, err := appFS.Open(m.Path)
	if err != nil {
		return nil, 0, err
	}
	return r, f.Size(), nil
}

// Name returns the filename.
func (m FileItem) Name() string {
	return path.Base(m.Path)
}

func (m FileItem) String() string {
	return m.Path
}

func (m FileItem) Size() int64 {
	f, err := appFS.Stat(m.Path)
	if err != nil {
		return 0
	}
	return f.Size()
}

func (m FileItem) Remove() error {
	return appFS.Remove(m.Path)
}
