// +build go1.7

package testhelper

import "testing"

func Run(t *testing.T, name string, f func(*testing.T)) {
	t.Run(name, f)
}
