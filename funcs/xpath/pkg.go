package xpath

import (
	"fmt"
	"io"

	"gopkg.in/xmlpath.v2"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

// Pkg contains custom functions defined by this package.
var Pkg = use.FuncMap{
	// func xpathStrict(path string, r io.Reader) (string, error)
	// Parses HTML, represented by r, compiles path and evaluates it to string.
	// Returns error, if cannot find node.
	// See "gopkg.in/xmlpath.v2".ParseHTML(), "gopkg.in/xmlpath.v2".Compile() and
	// "gopkg.in/xmlpath.v2".String().
	"xpathStrict": xpathStrict,
	// func xpath(path string, r io.Reader) (string, error)
	// In contrast to xpathStrict() silently return empty string, if cannot
	// find node.
	"xpath": xpathLoose,
}

func xpath(path string, r io.Reader) (string, bool, error) {
	node, err := xmlpath.ParseHTML(r)
	if err != nil {
		return "", false, err
	}
	p, err := xmlpath.Compile(path)
	if err != nil {
		return "", false, err
	}
	s, ok := p.String(node)
	return s, ok, nil
}

func xpathStrict(path string, r io.Reader) (string, error) {
	s, ok, err := xpath(path, r)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("xpath: path does not evaluate to string: %s", path)
	}
	return s, nil
}

func xpathLoose(path string, r io.Reader) (string, error) {
	s, _, err := xpath(path, r)
	if err != nil {
		return "", err
	}
	return s, nil
}
