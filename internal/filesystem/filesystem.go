package filesystem

import (
	"os/user"
	"path/filepath"
	"strings"
)

// AbsolutePath converts a path (relative or absolute) into an absolute one.
// Supports '~' notation for $HOME directory of the current user.
func AbsolutePath(path string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	dir := usr.HomeDir

	// In case of "~", which won't be caught by the next case
	if path == "~" {
		return dir, nil
	}

	// Use strings.HasPrefix so we don't match paths like
	// "/something/~/something/"
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(dir, path[2:]), nil
	}
	return filepath.Abs(path)
}
