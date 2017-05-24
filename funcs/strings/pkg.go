package strings

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

// Shifted (in contrast with "strings" package) args are more convenient with
// pipes int "text/template".

// Pkg contains custom functions defined by this package.
var Pkg = use.FuncMap{
	"atoi":   strconv.Atoi,
	"fields": strings.Fields,
	"lower":  strings.ToLower,
	"reader": strings.NewReader,
	// func replace(old, new, s string) string
	// Like string.Replace, but with reordered arguments.
	"replace": replace,
	// func split(sep, s string) []string
	// Like string.Split, but with reordered arguments.
	"split": split,
	// func string(v interface{}) string
	// Attempts to convert argument to string using reflection.
	"string":    convert,
	"trimSpace": strings.TrimSpace,
	"upper":     strings.ToUpper,
	// func contains(substr, str) bool
	// Like string.Contains, but with reordered arguments.
	"contains": contains,
}

func convert(v interface{}) string {
	return reflect.ValueOf(v).Convert(reflect.TypeOf("")).String()
}

func replace(old, new, s string) string {
	return strings.Replace(s, old, new, -1)
}

func split(sep, s string) []string {
	return strings.Split(s, sep)
}

func contains(substr, str string) bool {
	return strings.Contains(str, substr)
}
