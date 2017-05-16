package structor

import (
	"bytes"
	"fmt"
	"text/template"
)

// DefaultInterpreter is a default implementation of ELInterpreter,
// which is based on `"text/template"`.
type DefaultInterpreter struct {
	customFuncs template.FuncMap
}

// Execute implements ELInterpreter.Execute()
func (i *DefaultInterpreter) Execute(
	expression string,
	ctx *Context) (interface{}, error) {
	customFuncs := template.FuncMap{}
	for k, v := range i.customFuncs {
		customFuncs[k] = v
	}
	var res interface{}
	customFuncs["set"] = func(r interface{}) interface{} {
		res = r
		return r
	}
	templName := fmt.Sprintf("<<%T.%s>>", ctx.Struct, ctx.Name)
	t, err := template.New(templName).Funcs(customFuncs).Parse(expression)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, ctx)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}
	return buf.String(), nil
}
