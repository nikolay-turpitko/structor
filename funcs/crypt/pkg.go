package crypt

import (
	"strings"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

var Pkg = use.FuncMap{
	"rot13": func(s string) string { return strings.Map(rot13, s) },
}

func rot13(r rune) rune {
	if r >= 'a' && r <= 'z' {
		if r >= 'm' {
			return r - 13
		}
		return r + 13
	} else if r >= 'A' && r <= 'Z' {
		if r >= 'M' {
			return r - 13
		}
		return r + 13
	}
	return r
}
