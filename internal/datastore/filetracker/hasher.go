package filetracker

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// SHA256Hasher computes a cryptographic digest suitable for backup verification.
type SHA256Hasher struct{}

// Hash returns the SHA-256 of the file specified by filename.
func (h SHA256Hasher) Hash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
