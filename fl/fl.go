package fl

import "os"

// IsExists returns true when file is exists.
func IsExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
