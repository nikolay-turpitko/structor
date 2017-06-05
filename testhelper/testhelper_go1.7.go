// +build go1.7

package testhelper

import "testing"

// Run is a helper function to invoke testing.T.Run() when it is available.
func Run(t *testing.T, name string, f func(*testing.T)) {
	t.Run(name, f)
}
