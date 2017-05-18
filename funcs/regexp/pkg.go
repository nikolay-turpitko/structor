package regexp

import (
	"regexp"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

var Pkg = use.FuncMap{
	"match": func(s, re string) [][]string {
		return regexp.MustCompile(re).FindAllStringSubmatch(s, -1)
	},
}
