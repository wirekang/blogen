package fl

import "os"

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
