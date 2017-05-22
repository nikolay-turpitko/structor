package bytes

import (
	"bytes"
	"reflect"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

// Pkg contains custom functions defined by this package.
//
// func bytes(interface{}) []byte
//
// Attempts to convert argument to []byte using reflection.
//
// func reader(b []byte) *bytes.Reader
//
// Creates io.Reader from []byte.
var Pkg = use.FuncMap{
	"bytes":  convert,
	"reader": bytes.NewReader,
}

func convert(v interface{}) []byte {
	return reflect.ValueOf(v).Convert(reflect.TypeOf([]byte{})).Bytes()
}
