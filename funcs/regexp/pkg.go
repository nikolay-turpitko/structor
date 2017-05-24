package regexp

import (
	"regexp"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

// Pkg contains custom functions defined by this package.
var Pkg = use.FuncMap{
	// func(re, s string) [][]string
	// Returns regexp.FindAllStringSubmatch.
	"match": match,
	// func indx(i,j int, a [][]string) string
	// Indexes array and silently returns empty string if indexes are out of range.
	"indx": indx,
}

func match(re, s string) ([][]string, error) {
	r, err := regexp.Compile(re)
	if err != nil {
		return nil, err
	}
	return r.FindAllStringSubmatch(s, -1), nil
}

func indx(i, j int, a [][]string) string {
	if i < len(a) && j < len(a[i]) {
		return a[i][j]
	}
	return ""
}
