// Package fl serves uitl functions about files.
package fl

import (
	"crypto/md5"
	"io"
	"os"

	"github.com/wirekang/blogen/er"
)

// IsExist returns true if file exist.
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// MakeIfNotExist makes directory if it doesn't exist.
func MakeIfNotExist(path string) error {
	if !IsExist(path) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// CreateIfNotExist creates file if it doesn't exist.
func CreateIfNotExist(name string) (*os.File, error) {
	if !IsExist(name) {
		return os.Create(name)
	}
	return nil, nil
}

// Checksum returns hash of file.
func Checksum(file string) ([]byte, error) {
	hasher := md5.New()
	f, err := os.Open(file)
	if er.PrintIfNotNil(err) {
		return nil, err
	}
	defer f.Close()
	_, err = io.Copy(hasher, f)
	return hasher.Sum(nil), err
}
