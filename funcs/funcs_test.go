package funcs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nikolay-turpitko/structor/funcs/math"
	"github.com/nikolay-turpitko/structor/funcs/strings"
	"github.com/nikolay-turpitko/structor/funcs/use"
)

func TestProvide(t *testing.T) {
	ff := use.Packages(
		use.Pkg{Prefix: "", Funcs: math.Pkg},
		use.Pkg{Prefix: "", Funcs: strings.Pkg},
		use.Pkg{Prefix: "my_", Funcs: use.FuncMap{
			"echo": func(s string) string {
				return s + " " + s
			},
		},
		},
	)

	assert.Contains(t, ff, "my_echo")
	assert.Contains(t, ff, "add")
	assert.Contains(t, ff, "split")
}
