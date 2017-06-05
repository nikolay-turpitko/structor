// +build !go1.7

package testhelper

import (
	"fmt"
	"testing"
)

// Run is a helper function to emulate testing.T.Run() when it is not available.
func Run(t *testing.T, name string, f func(*testing.T)) {
	fmt.Printf("### Running: %s\n", name)
	f(t)
}
