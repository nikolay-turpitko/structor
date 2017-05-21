package strings

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

// Shifted (in contrast with "strings" package) args are more convenient with
// pipes int "text/template".

var Pkg = use.FuncMap{
	"atoi":   strconv.Atoi,
	"fields": strings.Fields,
	"split": func(sep, s string) []string {
		return strings.Split(s, sep)
	},
	"trimSpace": strings.TrimSpace,
	"upper":     strings.ToUpper,
	"lower":     strings.ToLower,
	"replace": func(old, new, s string) string {
		return strings.Replace(s, old, new, -1)
	},
	"string": func(v interface{}) string {
		return reflect.ValueOf(v).Convert(reflect.TypeOf("")).String()
	},
	"reader": strings.NewReader,
}
