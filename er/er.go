// Package er handles error
package er

import "fmt"

// PrintIfNotNil prints error and returns true if error isn't nil.
func PrintIfNotNil(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}
