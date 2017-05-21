package goquery

import (
	"io"

	"github.com/PuerkitoBio/goquery"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

var Pkg = use.FuncMap{
	"goquery": func(selector string, r io.Reader) (*goquery.Selection, error) {
		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			return nil, err
		}
		return doc.Find(selector), nil
	},
}
