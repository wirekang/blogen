package fl

import "os"

// IsExists returns true if file exist.
func IsExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// MakeIfNotExists makes directory if it doesn't exist.
func MakeIfNotExists(path string) error {
	if !IsExists(path) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}
