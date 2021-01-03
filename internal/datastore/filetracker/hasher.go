package filetracker

import (
	"fmt"
	"io"
	"os"

	"github.com/pierrec/xxHash/xxHash32"
)

// xxHash32Hasher implements a Hasher using xxHash32 package.
type xxHash32Hasher struct{}

// Hash returns the xxHash32 of the file specified by filename.
func (h xxHash32Hasher) Hash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	hasher := xxHash32.New(0xCAFE)
	defer hasher.Reset()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(hasher.Sum32()), nil
}
