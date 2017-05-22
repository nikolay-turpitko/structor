package goquery

import (
	"io"

	"github.com/PuerkitoBio/goquery"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

// Pkg contains custom functions defined by this package.
var Pkg = use.FuncMap{
	// func goquery(selector string, r io.Reader) (*goquery.Selection, error)
	// Finds selector in the HTML document, represented by r.
	// See "github.com/PuerkitoBio/goquery".NewDocumentFromReader() and
	// "github.com/PuerkitoBio/goquery".Find().
	"goquery": goQuery,
}

func goQuery(selector string, r io.Reader) (*goquery.Selection, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	return doc.Find(selector), nil
}
