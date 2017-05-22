package encoding

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

// Pkg contains custom functions defined by this package.
var Pkg = use.FuncMap{
	"base64":   base64.StdEncoding.EncodeToString,
	"unbase64": base64.StdEncoding.DecodeString,
	"hex":      hex.EncodeToString,
	"unhex":    hex.DecodeString,
}
