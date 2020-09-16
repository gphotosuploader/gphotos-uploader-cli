package upload

import (
	"io"
	"os"
	"path"
)

// FileItem represents a local file.
type FileItem struct {
	Path string
	AlbumName string
}

// Open returns a stream.
// Caller should close it finally.
func (m FileItem) Open() (io.ReadSeeker, int64, error) {
	f, err := os.Stat(m.String())
	if err != nil {
		return nil, 0, err
	}
	r, err := os.Open(m.String())
	if err != nil {
		return nil, 0, err
	}
	return r, f.Size(), nil
}

// Name returns the filename.
func (m FileItem) Name() string {
	return path.Base(m.String())
}

func (m FileItem) String() string {
	return string(m.Path)
}

func (m FileItem) Size() int64 {
	f, err := os.Stat(m.String())
	if err != nil {
		return 0
	}
	return f.Size()
}
