package encoding

import (
	"encoding/base64"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

var Pkg = use.FuncMap{
	"unbase64": base64.StdEncoding.DecodeString,
}
