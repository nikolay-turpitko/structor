package bytes

import (
	"bytes"
	"reflect"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

var Pkg = use.FuncMap{
	"bytes": func(v interface{}) []byte {
		return reflect.ValueOf(v).Convert(reflect.TypeOf([]byte{})).Bytes()
	},
	"reader": bytes.NewReader,
}
