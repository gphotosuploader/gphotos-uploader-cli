package filetracker

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// TrackedFile represents a tracked file in the repository.
type TrackedFile struct {
	Version     int       `json:"version,omitempty"`
	ModTime     time.Time `json:"mod_time"`
	Size        int64     `json:"size,omitempty"`
	Hash        string    `json:"sha256"`
	MediaItemID string    `json:"media_item_id,omitempty"`
	ProductURL  string    `json:"product_url,omitempty"`
	UploadedAt  time.Time `json:"uploaded_at,omitempty"`
}

// NewTrackedFile returns a TrackedFile with the specified values
func NewTrackedFile(value string) TrackedFile {
	if strings.HasPrefix(value, "{") {
		var tracked TrackedFile
		if json.Unmarshal([]byte(value), &tracked) == nil {
			return tracked
		}
	}
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
		Hash:    hash,
		ModTime: modTime,
	}
}

func (tf TrackedFile) String() string {
	if tf.Version >= 2 {
		encoded, _ := json.Marshal(tf)
		return string(encoded)
	}
	if tf.ModTime.IsZero() {
		return tf.Hash
	} else {
		return strconv.FormatInt(tf.ModTime.UnixNano(), 10) + "|" + tf.Hash
	}
}
