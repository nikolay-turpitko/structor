package structor

import (
	"encoding/base64"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

// providedFuncs - map of functions for default EL interpreter.
var providedFuncs = template.FuncMap{
	"fields": strings.Fields,
	"match": func(s, re string, i, j int) string {
		return regexp.MustCompile(re).FindAllStringSubmatch(s, -1)[i][j]
	},
	"split":     strings.Split,
	"trimSpace": strings.TrimSpace,
	"atoi":      strconv.Atoi,
	"unbase64":  base64.StdEncoding.DecodeString,
}

// AddProvidedFuncs copies providedFuncs and customFuncs into the new map.
// Note: custom functions with the same name will replace provided ones (by
// design).
func AddProvidedFuncs(customFuncs template.FuncMap) template.FuncMap {
	funcs := template.FuncMap{}
	for k, v := range providedFuncs {
		funcs[k] = v
	}
	if customFuncs != nil {
		for k, v := range customFuncs {
			funcs[k] = v
		}
	}
	return funcs
}
