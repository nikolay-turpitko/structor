package strings

import (
	"strconv"
	"strings"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

var Pkg = use.FuncMap{
	"atoi":      strconv.Atoi,
	"fields":    strings.Fields,
	"split":     strings.Split,
	"trimSpace": strings.TrimSpace,
}
