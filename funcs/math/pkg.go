package math

import (
	"reflect"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

type oprnd interface{}
type operation func(float64, float64) float64

// Pkg contains custom functions defined by this package.
var Pkg = use.FuncMap{
	// Adds all operands, converting them to float64.
	"add": add,
	// Subtracts all operands from the first, converting them to float64.
	"sub": sub,
	// Multiplies all operands, converting them to float64.
	"mul": mul,
	// Divides first operand on all other operands, converting them to float64.
	"div": div,
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

func add(op ...oprnd) oprnd {
	return perform(func(a, b float64) float64 { return a + b }, op...)
}
func sub(op ...oprnd) oprnd {
	return perform(func(a, b float64) float64 { return a - b }, op...)
}
func mul(op ...oprnd) oprnd {
	return perform(func(a, b float64) float64 { return a * b }, op...)
}
func div(op ...oprnd) oprnd {
	return perform(func(a, b float64) float64 { return a / b }, op...)
}
