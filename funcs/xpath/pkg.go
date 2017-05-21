package xpath

import (
	"fmt"
	"io"

	"gopkg.in/xmlpath.v2"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

var Pkg = use.FuncMap{
	"xpath": func(path string, r io.Reader) (string, error) {
		node, err := xmlpath.ParseHTML(r)
		if err != nil {
			return "", err
		}
		p, err := xmlpath.Compile(path)
		if err != nil {
			return "", err
		}
		s, ok := p.String(node)
		if !ok {
			return "", fmt.Errorf("xpath: path does not evaluate to string: %s", path)
		}
		return s, nil
	},
}
