package xpath

import (
	"fmt"
	"io"

	"gopkg.in/xmlpath.v2"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

// Pkg contains custom functions defined by this package.
var Pkg = use.FuncMap{
	// func xpath(path string, r io.Reader) (string, error)
	// Parses HTML, represented by r, compiles path and evaluates it to string.
	// See "gopkg.in/xmlpath.v2".ParseHTML(), "gopkg.in/xmlpath.v2".Compile() and
	// "gopkg.in/xmlpath.v2".String().
	"xpath": xpath,
}

func xpath(path string, r io.Reader) (string, error) {
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
}
