package bytes

import (
	"bytes"
	"reflect"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

// Pkg contains custom functions defined by this package.
var Pkg = use.FuncMap{
	// Signature: func bytes(interface{}) []byte.
	// Attempts to convert argument to []byte using reflection.
	"bytes": convert,
	// Signature: func reader(b []byte) *bytes.Reader.
	// Creates io.Reader from []byte.
	"reader": bytes.NewReader,
}

func convert(v interface{}) []byte {
	return reflect.ValueOf(v).Convert(reflect.TypeOf([]byte{})).Bytes()
}
