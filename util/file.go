package util

import (
	"fmt"
	"os"
	"path"
)

// MaybeCreateFile creates file if not exists
func MaybeCreateFile(dir, subDir, file string, cover bool) (fp *os.File, created bool, err error) {
	err = MkdirIfNotExist(path.Join(dir, subDir))
	if err != nil {
		return nil, false, err
	}
	filePath := path.Join(dir, subDir, file)
	if FileExists(filePath) {
		if cover {
			fmt.Printf("%s exists, covered generation\n", filePath)
			err := RemoveIfExist(filePath)
			if err != nil {
				return nil, false, err
			}
		} else {
			fmt.Printf("%s exists, ignored generation\n", filePath)
			return nil, false, nil
		}
	}
	fp, err = CreateIfNotExist(filePath)
	created = err == nil
	return
}

// MkdirIfNotExist makes directories if the input path is not exists
func MkdirIfNotExist(dir string) error {
	if len(dir) == 0 {
		return nil
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}

	return nil
}

// FileExists returns true if the specified file is exists.
func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// CreateIfNotExist creates a file if it is not exists.
func CreateIfNotExist(file string) (*os.File, error) {
	_, err := os.Stat(file)
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s already exist", file)
	}

	return os.Create(file)
}

// RemoveIfExist deletes the specified file if it is exists.
func RemoveIfExist(filename string) error {
	if !FileExists(filename) {
		return nil
	}

	return os.Remove(filename)
}
