package use

import (
	"log"
	"text/template"
)

// FuncMap is a template.FuncMap, see docs there. Redefined to minimize imports.
type FuncMap template.FuncMap

// Pkg described single "package" of related functions.
// "Package" can be included into the map with every function, prefixed with
// Prefix. This is namespaces mechanism for poor.
// MapNames is a function invoking for each key within Funcs to convert names
// using custom rules. This can be used to adapt names to use with different
// interpreters.
type Pkg struct {
	Prefix  string
	MapName func(string) string
	Funcs   FuncMap
}

// Packages collects functions from all "packages" in arguments into one
// FuncMap, converting every function name in "package" with Pkg.MapName and
// prefixing it with Pkg.Prefix.
func Packages(pkgs ...Pkg) FuncMap {
	l := 0
	for _, p := range pkgs {
		l += len(p.Funcs)
	}
	m := make(FuncMap, l)
	for _, p := range pkgs {
		for nm, f := range p.Funcs {
			if p.MapName != nil {
				nm = p.MapName(nm)
			}
			name := p.Prefix + nm
			if f2, ok := m[name]; ok {
				log.Printf("use: name clash for '%s', %T, %T", name, f, f2)
			}
			m[name] = f
		}
	}
	return m
}
