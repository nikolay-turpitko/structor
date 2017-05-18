package math

import (
	"reflect"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

type oprnd interface{}
type operation func(float64, float64) float64

var Pkg = use.FuncMap{
	"add": func(op ...oprnd) oprnd {
		return perform(func(a, b float64) float64 { return a + b }, op...)
	},
	"sub": func(op ...oprnd) oprnd {
		return perform(func(a, b float64) float64 { return a - b }, op...)
	},
	"mul": func(op ...oprnd) oprnd {
		return perform(func(a, b float64) float64 { return a * b }, op...)
	},
	"div": func(op ...oprnd) oprnd {
		return perform(func(a, b float64) float64 { return a / b }, op...)
	},
}

func toIntrnl(op oprnd) float64 {
	return reflect.ValueOf(op).Convert(reflect.TypeOf(float64(0))).Float()
}

func perform(f operation, op ...oprnd) oprnd {
	res := toIntrnl(op[0])
	for i, l := 1, len(op); i < l; i++ {
		res = f(res, toIntrnl(op[i]))
	}
	return res
}
