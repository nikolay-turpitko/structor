package regexp

import (
	"regexp"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

// Pkg contains custom functions defined by this package.
var Pkg = use.FuncMap{
	"match": func(re, s string) [][]string {
		return regexp.MustCompile(re).FindAllStringSubmatch(s, -1)
	},
}
