// +build !go1.7

package testhelper

import (
	"fmt"
	"testing"
)

func Run(t *testing.T, name string, f func(*testing.T)) {
	fmt.Printf("### Running: %s\n", name)
	f(t)
}
