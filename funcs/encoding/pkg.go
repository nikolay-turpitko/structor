package encoding

import (
	"encoding/base64"
	"encoding/hex"
	"reflect"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

var Pkg = use.FuncMap{
	"base64":   base64.StdEncoding.EncodeToString,
	"unbase64": base64.StdEncoding.DecodeString,
	"hex":      hex.EncodeToString,
	"unhex":    hex.DecodeString,
	"bytes": func(v interface{}) []byte {
		return reflect.ValueOf(v).Convert(reflect.TypeOf([]byte{})).Bytes()
	},
}
