// +build go1.7

package goel

import (
	"fmt"
	"reflect"

	"github.com/apaxa-go/eval"

	"github.com/nikolay-turpitko/structor/el"
)

// Interpreter implements "github.com/nikolay-turpitko/structor/el".Interpreter
// using "github.com/apaxa-go/eval". Thus, it evaluates Go-style expressions.
//
// Besides eval.Args passed to it, Interpreter creates these additional
// custom objects for use in EL expressions:
//
//  - ctx
//  - ctxStruct
//  - ctxExtra
//  - ctxSub
//  - eval
//
// Structure "ctx" is a context of type *el.Context.
//
// Types "ctxStruct", "ctxExtra" and "ctxSub" are the EL aliases of actual
// types of correspondent fields of "ctx" and can be used to convert inteface{}
// types of these fields to proper types to access their fields.
//
// Function "eval" with signature `func(intrpr, expr string) interface{}`
// can be used to evaluate given expression with given interpreter.
// This can be useful to evaluate expression passed as a string within context.
// This feature works with support of the calling code, which passes actual
// implementation of the "eval" function in Context.EvalExpr().
// Interpreter name should be known to calling code. For example, for known
// implementation (structor.NewEvaluator()) interpreter name is a tag name,
// onto which given interpreter is mapped during creation of evaluator.
//
// Due restrictions of "github.com/apaxa-go/eval", only custom functions
// returning one or two results are  allowed. If custom function returns two
// results, its second result must be of error type and it's converted to
// panic (which is catched and converted back to error by the
// "github.com/apaxa-go/eval").
type Interpreter struct {
	// Arguments for expression.
	Args eval.Args
}

// Execute implements Interpreter.Execute()
func (i *Interpreter) Execute(
	expression string,
	ctx *el.Context) (interface{}, error) {
	expr, err := eval.ParseString(expression, "")
	if err != nil {
		return nil, fmt.Errorf("structor parse: <<%s>>: %v", ctx.LongName, err)
	}
	funcEval := func(intrpr, expr string) interface{} {
		res, err := ctx.EvalExpr(intrpr, expr, ctx)
		if err != nil {
			panic(err)
		}
		return res
	}
	args := eval.Args{
		"eval":      eval.MakeDataRegularInterface(funcEval),
		"ctx":       eval.MakeDataRegularInterface(ctx),
		"ctxStruct": eval.MakeTypeInterface(ctx.Struct),
	}
	if ctx.Extra != nil {
		args["ctxExtra"] = eval.MakeTypeInterface(ctx.Extra)
	}
	if ctx.Sub != nil {
		args["ctxSub"] = eval.MakeTypeInterface(ctx.Sub)
	}
	for k, v := range i.Args {
		args[k] = wrapFunc(v)
	}
	res, err := expr.EvalToInterface(args)
	if err != nil {
		return nil, fmt.Errorf("structor eval: <<%s>>: %v", ctx.LongName, err)
	}
	return res, nil
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// wrapFunc check if argument is function with two return values, last of which
// is error, and wraps such a function to return only one value, as apaxa-go
// permits.
func wrapFunc(v eval.Value) eval.Value {
	if v.Kind() != eval.Datas {
		return v
	}
	if v.Data().Kind() != eval.Regular {
		return v
	}
	r := v.Data().Regular()
	if r.Kind() != reflect.Func {
		return v
	}
	t := r.Type()
	n := t.NumOut()
	if n != 2 {
		return v
	}
	if t.Out(1) != errorType {
		return v
	}
	in := make([]reflect.Type, 0, t.NumIn())
	for i, l := 0, t.NumIn(); i < l; i++ {
		in = append(in, t.In(i))
	}
	out := []reflect.Type{t.Out(0)}
	twraper := reflect.FuncOf(in, out, false)
	wraper := reflect.MakeFunc(twraper, func(args []reflect.Value) []reflect.Value {
		result := r.Call(args)
		if len(result) == 2 && !result[1].IsNil() {
			panic(result[1])
		}
		return result[:1]
	})
	return eval.MakeDataRegular(wraper)
}
