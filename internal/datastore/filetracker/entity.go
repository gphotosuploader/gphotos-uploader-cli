package filetracker

import (
	"strings"
	"time"
	"strconv"
)

// TrackedFile represents a tracked file in the repository.
type TrackedFile struct {
	ModTime time.Time
	Hash string
}

// NewTrackedFile returns a TrackedFile with the specified values
func NewTrackedFile(value string) TrackedFile {
	parts := strings.SplitN(value, "|", 2)

	modTime := time.Time{}
	hash := ""

	if len(parts) == 2 {
		unixTime, err := strconv.ParseInt(parts[0], 10, 64)
		if err == nil {
			modTime = time.Unix(0, unixTime)
		}
		hash = parts[1]
	} else {
		hash = parts[0]
	}

	return TrackedFile{
		Hash: hash,
		ModTime: modTime,
	}
}

func (tf TrackedFile) String() string {
	if tf.ModTime.IsZero() {
		return tf.Hash
	} else {
		return strconv.FormatInt(tf.ModTime.UnixNano(), 10) + "|" + tf.Hash
	}
}