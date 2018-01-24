/*
Package el provides an interface and default implementation of expression
language (EL) interpreter for struct tags.

Default implementation is simply based on "text/template".
*/
package el

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

// Interpreter is an interface of EL interpreter.
type Interpreter interface {
	// Execute parses and executes expression with a given context.
	Execute(expression string, ctx *Context) (result interface{}, err error)
}

// The InterpreterFunc type is an adapter to allow the use of ordinary
// functions as EL interpreter. If f is a function with the appropriate
// signature, InterpreterFunc(f) is a Interpreter that calls f.
type InterpreterFunc func(string, *Context) (interface{}, error)

// Execute calls f(s, ctx).
func (f InterpreterFunc) Execute(s string, ctx *Context) (interface{}, error) {
	return f(s, ctx)
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
	// Function, which knows how to evaluate expression with different
	// interpreter.
	EvalExpr EvalExprFunc
}

// EvalExprFunc is a type of function, which knows how to evaluate given
// expression using given interpreter name and context.
// It is used to implement special predefined custom function "eval", available
// within EL.
type EvalExprFunc func(
	interpreterName, expression string,
	context *Context) (interface{}, error)

// DefaultInterpreter is a default implementation of Interpreter,
// which is based on "text/template".
//
// Besides Funcs passed to it, DefaultInterpreter creates these additional
// custom function for use in EL expressions:
//
//  - set
//  - eval
//
// Function "set" with signature `func (r interface{}) interface{}` passes
// argument to result, but stores it internally to be used as an expression
// result. It's intention is to set expression result to some concrete type,
// other than string (which is only type "text/template" allows by design).
// Simply put, when you interested in interface{} result and not in it's
// default string representation, use "set" to store result at the end of
// expression.
//
// Function "eval" with signature `func(intrpr, expr string) (interface{}, error)`
// can be used to evaluate given expression with given interpreter.
// This can be useful to evaluate expression passed as a string within context.
// This feature works with support of the calling code, which passes actual
// implementation of the "eval" function in Context.EvalExpr().
// Interpreter name should be known to calling code. For example, for known
// implementation (structor.NewEvaluator()) interpreter name is a tag name,
// onto which given interpreter is mapped during creation of evaluator.
//
// Restrictions of "text/template" package applied to custom functions.
type DefaultInterpreter struct {
	// Custom functions, available for use in EL expressions.
	Funcs use.FuncMap
	// Left delimiter for templates.
	LeftDelim string
	// Right delimiter for templates.
	RightDelim string
	// Automatically enclose passed expression into delimiters before
	// interpretation (if it is not already enclosed). This allows to pass
	// simplified expressions. For example, `atoi "42"` instead of
	// `{{atoi "42"}}`.
	AutoEnclose bool
}

// Execute implements Interpreter.Execute()
func (i *DefaultInterpreter) Execute(
	expression string,
	ctx *Context) (interface{}, error) {
	funcs := template.FuncMap{}
	for k, v := range i.Funcs {
		funcs[k] = v
	}
	var res interface{}
	resultEvaluated := false
	funcs["set"] = func(r interface{}) interface{} {
		res = r
		resultEvaluated = true
		return r
	}
	funcs["eval"] = func(intrpr, expr string) (interface{}, error) {
		return ctx.EvalExpr(intrpr, expr, ctx)
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
	if i.AutoEnclose &&
		!(strings.HasPrefix(expression, left) &&
			strings.HasSuffix(expression, right)) {
		expression = fmt.Sprintf("%s%s%s", left, expression, right)
	}
	t, err := template.
		New(templName).
		Delims(left, right).
		Funcs(funcs).
		Parse(expression)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, ctx)
	if err != nil {
		return nil, err
	}
	if resultEvaluated {
		return res, nil
	}
	return buf.String(), nil
}
