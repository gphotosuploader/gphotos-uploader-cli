package upload

import (
	"io"
	"path"

	"github.com/spf13/afero"
)

var (
	// fs represents the filesystem to use. By default, it uses functions based on the `os` package.
	// In testing, it uses a memory file system.
	appFS = afero.NewOsFs()
)

// FileItem represents a local file.
type FileItem struct {
	Path      string
	AlbumName string
}

// NewFileItem creates a new instance of FileItem.
func NewFileItem(path string) FileItem {
	return FileItem{
		Path: path,
	}
}

// Open opens the file and returns a stream.
// The caller should close it finally.
func (f FileItem) Open() (io.ReadSeeker, int64, error) {
	fileInfo, err := appFS.Stat(f.Path)
	if err != nil {
		return nil, 0, err
	}

	file, err := appFS.Open(f.Path)
	if err != nil {
		return nil, 0, err
	}

	return file, fileInfo.Size(), nil
}

// Name returns the filename.
func (f FileItem) Name() string {
	return path.Base(f.Path)
}

// String returns the string representation of the FileItem.
func (f FileItem) String() string {
	return f.Path
}

// Size returns the file size.
func (f FileItem) Size() int64 {
	fileInfo, err := appFS.Stat(f.Path)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

// Remove removes the file from the file system.
func (f FileItem) Remove() error {
	return appFS.Remove(f.Path)
}

// GroupByAlbum groups FileItem objects by their AlbumName.
func GroupByAlbum(items []FileItem) map[string][]FileItem {
	groups := make(map[string][]FileItem)

	for _, item := range items {
		groups[item.AlbumName] = append(groups[item.AlbumName], item)
	}

	return groups
}
