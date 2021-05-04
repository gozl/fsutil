package fsutil

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	ErrNotFile = errors.New("not a file")
	ErrFileTooBig = errors.New("file is too large")
)

// RemoveFile deletes file at path. It returns an error if path is not a file. 
// Set deleteEmptyDir to true will also delete the file's parent directory if 
// that parent directory will be empty after deleting file.
func RemoveFile(path string, deleteEmptyDir bool) error {
	if TestPath(path) != File {
		return ErrNotFile
	}

	err := os.Remove(path)
	if err != nil {
		return err
	}

	if !deleteEmptyDir {
		return nil
	}

	parentDir := filepath.Dir(path)
	ok, err := IsEmptyDir(parentDir)
	if err != nil {
		return err
	}
	if ok {
		return os.Remove(parentDir)
	}
	return nil
}

// ReadFile reads all content of file at path into memory. If maxBytes is 
// larger than 0, returns (nil,error) if the file size is larger than maxBytes. 
// Returns (nil,error) if path is not a file.
func ReadFile(path string, maxBytes int64) ([]byte, error) {
	if maxBytes > 0 {
		stat, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if stat.IsDir() {
			return nil, ErrNotFile
		}
		if maxBytes > 0 && stat.Size() > maxBytes {
			return nil, ErrFileTooBig
		}
	}
	return os.ReadFile(path)
}

// WriteFile overwrites content of file at path with data, and then sets the 
// file permission to mode. The file at path is created automatically if it 
// does not exist, but an error occurs if the file's parent directory does not 
// exist.
// If a file at path already exists, it is overwritten. Use AppendFile if you 
// wish to append content instead.
func WriteFile(path string, data []byte, mode fs.FileMode) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.Write(data)
	if err != nil {
		return err
	}

	err = os.Chmod(path, mode)
	if err != nil {
		return err
	}
	return nil
}

// AppendFile appends data to file at path. If the file does not exist, it is 
// created with permission set to mode. If the file already exists, it mode is 
// ignored.
// An error occurs if the file's parent directory does not exist.
func AppendFile(path string, data []byte, mode fs.FileMode) error {
	fd, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.Write(data)
	if err != nil {
		return err
	}
	return nil
}
