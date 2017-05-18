/*
Package el provides an interface and default implementation of expression
language (EL) interpreter for struct tags.

Default implementation is simply based on "text/template".
*/
package el

import (
	"bytes"
	"fmt"
	"text/template"
)

// Interpreter is an interface of EL interpreter.
type Interpreter interface {
	Execute(expression string, ctx *Context) (result interface{}, err error)
}

// Context is a context, passed to interpreter.
// It contains information about currently processed field, struct and extra.
type Context struct {
	// Name of the currently processed field.
	Name string
	// Name of the currently processed field including type.
	LongName string
	// Current value of the currently processed filed.
	Val interface{}
	// All other tags of the currently processed field.
	Tags map[string]string
	// Currently processed struct.
	Struct interface{}
	// Extra context structure.
	Extra interface{}
	// Temporary partial result evaluated on the current substruct.
	Sub interface{}
}

// DefaultInterpreter is a default implementation of Interpreter,
// which is based on "text/template".
type DefaultInterpreter struct {
	// Custom functions, available for use in EL expressions.
	CustomFuncs template.FuncMap
	// Left delimiter for templates.
	LeftDelim string
	// Right delimiter for templates.
	RightDelim string
	// Automatically enclose passed expression into delimiters before
	// interpretation. This allows to pass simplified expressions. For example,
	// `atoi "42"` instead of `{{atoi "42"}}`.
	AutoEnclose bool
}

// Execute implements Interpreter.Execute()
func (i *DefaultInterpreter) Execute(
	expression string,
	ctx *Context) (interface{}, error) {
	customFuncs := template.FuncMap{}
	for k, v := range i.CustomFuncs {
		customFuncs[k] = v
	}
	var res interface{}
	customFuncs["set"] = func(r interface{}) interface{} {
		res = r
		return r
	}
	templName := fmt.Sprintf("<<%s>>", ctx.LongName)
	left := i.LeftDelim
	right := i.RightDelim
	if left == "" {
		left = "{{"
	}
	if right == "" {
		right = "}}"
	}
	if i.AutoEnclose {
		expression = fmt.Sprintf("%s%s%s", left, expression, right)
	}
	t, err := template.
		New(templName).
		Delims(left, right).
		Funcs(customFuncs).
		Parse(expression)
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
