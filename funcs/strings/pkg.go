package strings

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

var Pkg = use.FuncMap{
	"atoi":      strconv.Atoi,
	"fields":    strings.Fields,
	"split":     strings.Split,
	"trimSpace": strings.TrimSpace,
	"replace": func(s, old, new string) string {
		return strings.Replace(s, old, new, -1)
	},
	"string": func(v interface{}) string {
		return reflect.ValueOf(v).Convert(reflect.TypeOf("")).String()
	},
}
