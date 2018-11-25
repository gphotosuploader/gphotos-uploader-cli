package fileshandling

import (
	"io/ioutil"
	"strings"

	"github.com/nmrshll/gphotos-uploader-cli/util"
	"github.com/palantir/stacktrace"
	filetype "gopkg.in/h2non/filetype.v1"
)

func IsFile(path string) bool {
	return util.IsFile(path)
}

func IsDir(path string) bool {
	return util.IsDir(path)
}

func IsImage(path string) (bool, error) {
	if !HasImageExtension(path) {
		return false, nil
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return false, stacktrace.Propagate(err, "Failed finding file type: %s: Ignoring file...\n", path)
	}

	kind, err := filetype.Match(buf)
	if err != nil {
		return false, stacktrace.Propagate(err, "Failed finding file type: %s: Ignoring file...\n", path)
	}

	if strings.Contains(kind.MIME.Value, "image") {
		return true, nil
	}
	return false, nil
}
