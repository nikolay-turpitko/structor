package math

import (
	"reflect"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

var Pkg = use.FuncMap{
	"add": func(a, b interface{}) interface{} { return convert(a) + convert(b) },
}

func convert(n interface{}) float64 {
	return reflect.ValueOf(n).Convert(reflect.TypeOf(float64(0))).Float()
}
