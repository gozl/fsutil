package fsutil

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrNotDir = errors.New("not a directory")
	ErrHomeDir = errors.New("cannot resolve home directory")
)

// PathType is the type of filesystem object.
type PathType uint32

const (
	// NotExist means the path does not exist.
	NotExist PathType = iota
	// Directory indicates a directory.
	Directory
	// File indicates a regular file.
	File
	// IrregularPath means the path exists but is an unsupported type.
	IrregularPath PathType = 54288
	// BadPath means the path exists, but more information about the path 
	// cannot be determined.
	BadPath PathType = 0655
)

// TestPath returns the path type.
func TestPath(path string) PathType {
	stat, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return NotExist
	}
	if err != nil {
		return BadPath
	}

	if stat.IsDir() {
		return Directory
	}

	filemode := stat.Mode()
	if filemode.IsRegular() {
		return File
	}
	return IrregularPath
}

// IsPath returns true if path is any of path type t.
func IsPath(path string, t ...PathType) bool {
	if len(t) == 0 {
		return true
	}

	pathType := TestPath(path)
	for _, v := range t {
		if v == pathType {
			return true
		}
	}
	return false
}

// IsEmptyDir returns true if the directory is empty. Returns ErrNotDir if path 
// is not a directory.
func IsEmptyDir(path string) (bool, error) {
	fd, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer fd.Close()

	_, err = fd.Readdirnames(1)
	if err != nil && err == io.EOF {
		return true, nil
	}
	if err != nil && strings.HasSuffix(err.Error(), "not a directory") {
		return false, ErrNotDir
	}
	return false, err
}

// RemoveEmptyDir removes the directory at path if it is empty. Returns 
// ErrNotDir if path is not a directory.
func RemoveEmptyDir(path string) error {
	ok, err := IsEmptyDir(path)
	if err != nil {
		return err
	}
	if ok {
		return os.Remove(path)
	}
	return nil
}

// Dir returns regular files directly in directory path up to maxfile entries.
// Use Subdir if you want to return sub-directories.
// If maxfile is <= 0, it returns all the entries.
// If ext is empty or *, entries are returned as-is. Otherwise, only entries
// with matching extension will be returned.
// Returns ErrNotDir if path is not a directory.
func Dir(path, ext string, maxfile int) ([]string, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	dirItems, err := fd.ReadDir(maxfile)
	if err != nil && err == io.EOF {
		return nil, nil
	}
	if err != nil && strings.HasSuffix(err.Error(), "not a directory") {
		return nil, ErrNotDir
	}

	// filepath.Ext(x) returns the dot too:
	// e.g. filepath.Ext("foo.txt") -> ".txt"
	fullExt := ""
	if ext != "*" && ext != "" {
		fullExt = "." + ext
	}

	names := []string{}
	for _, v := range dirItems {
		if !v.Type().IsRegular() {
			continue
		}
		if fullExt == "" {
			names = append(names, v.Name())
			continue
		}
		if filepath.Ext(v.Name()) == fullExt {
			names = append(names, v.Name())
		}
	}

	// keep in sync with os.Dir: don't return nil slices!
	return names, nil
}

// Subdir returns sub-directories directly in directory path up to maxfile 
// entries. If maxfile is <= 0, it returns all the entries.
// Use Dir if you want to return regular files.
// Returns ErrNotDir if path is not a directory.
func Subdir(path string, maxfile int) ([]string, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	dirItems, err := fd.ReadDir(maxfile)
	if err != nil && err == io.EOF {
		return nil, nil
	}
	if err != nil && strings.HasSuffix(err.Error(), "not a directory") {
		return nil, ErrNotDir
	}

	names := []string{}
	for _, v := range dirItems {
		if !v.IsDir() {
			continue
		}
		names = append(names, v.Name())
	}

	// keep in sync with os.Dir: don't return nil slices!
	return names, nil
}

// HomeDir returns the current user's home directory. It returns an empty string 
// if the user's home directory cannot be determined.
func HomeDir() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return homePath
}

func Abs(path string) (string, error) {
	var result string
	if strings.HasPrefix(path, "~/") {
		homedir := HomeDir()
		if homedir == "" {
			return "", ErrHomeDir
		}
		if len(path) == 2 {
			result = homedir
		} else {
			result = homedir + "/" + path[2:]
		}
	} else {
		result = path
	}

	return filepath.Abs(result)
}
