package use

import "text/template"

// FuncMap is a template.FuncMap, see docs there. Redefined to minimize imports.
type FuncMap template.FuncMap

// Pkg described single "package" of related functions.
// "Package" can be included into the map with every function, prefixed with
// Prefix. This is namespaces mechanism for poor.
type Pkg struct {
	Prefix string
	Funcs  FuncMap
}

// Packages collects functions from all "packages" in arguments into one
// FuncMap, prefixing every function name in "package" with Pkg.Prefix.
func Packages(pkgs ...Pkg) FuncMap {
	l := 0
	for _, p := range pkgs {
		l += len(p.Funcs)
	}
	m := make(FuncMap, l)
	for _, p := range pkgs {
		for nm, f := range p.Funcs {
			m[p.Prefix+nm] = f
		}
	}
	return m
}
