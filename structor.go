/*
Package structor contains interface and default implementation of EL
(expression language) evaluator which evaluates EL expressions within tags of
the struct, using optional additional struct as an extra context.

Basic idea is to use simple expression language within Go struct tags to
compute struct fields based on other fields or provided additional context.

Due usage of reflection and EL interpretation, this package is hardly suitable
for tasks, requiring high performance, but rather intended to be used during
application setup or in cases where high performance is not an ultimate goal.

See tests for examples of usage with xpath, regexp, goquery, encryption,
reading from files, shell invocation, etc.
*/
package structor

import (
	"fmt"
	"reflect"

	multierror "github.com/hashicorp/go-multierror"

	"github.com/nikolay-turpitko/structor/el"
	"github.com/nikolay-turpitko/structor/funcs/use"
	"github.com/nikolay-turpitko/structor/scanner"
)

// Evaluator is an interface of evaluator, which gets structure and extra
// context as input, iterates over `s`'s fields and evaluate expression tag on
// every field.
type Evaluator interface {
	Eval(s, extra interface{}) error
}

// Interpreters is a map of tag names to el.Interpreters.  Used to register
// different interpreters for different tag names.
//
// Only first tag name on the struct field is currently recognized and
// processed. So, only one EL expression per structure field, but different
// fields of the same structure can be processed by different interpreters.
type Interpreters map[string]el.Interpreter

// WholeTag constant can be used as tag name in the Interpreters to indicate
// that whole tag value should be passed to the interpreter.
//
// Registering interpreter to the whole tag value conflicts with any other
// usage of the struct's tag, but can be convenient to simplify complex EL
// expressions with quotes (regexp, for example).
//
// WholeTag interpreter is probed after all other registered interpreters.
const WholeTag = ""

// NewEvaluator returns Evaluator with desired settings.
//
// Only first tag with EL will be recognized and used (only one
// expression per struct field). Different fields of the same struct can be
// processed using different EL interpreters.
//
//  scanner - is a scanner implementation to be used to scan tags.
//  interpreters - is a map of registered tag names to EL interpreters.
func NewEvaluator(
	scanner scanner.Scanner,
	interpreters Interpreters) Evaluator {
	return NewEvaluatorWithOptions(scanner, interpreters, Options{})
}

func NewEvaluatorWithOptions(
	scanner scanner.Scanner,
	interpreters Interpreters,
	options Options) Evaluator {
	if len(interpreters) == 0 {
		panic("no interpreters registered")
	}
	return &evaluator{scanner, interpreters, options}
}

// NewDefaultEvaluator returns default Evaluator implementation. Default
// implementation uses tag "eval" for expressions and EL interpreter, based on
// `"text/template"`.
//
//  funcs - custom functions, available for interpreter;
func NewDefaultEvaluator(funcs use.FuncMap) Evaluator {
	return NewEvaluator(
		scanner.Default,
		Interpreters{
			"eval": &el.DefaultInterpreter{Funcs: funcs},
		})
}

// NewNonmutatingEvaluator creates Evaluator implementation which does not
// change original structure (does not save evaluated results) itself.
// Though, interpreters can change structures' fields as a side effect.
//
// See NewEvaluator() for additional information.
func NewNonmutatingEvaluator(
	scanner scanner.Scanner,
	interpreters Interpreters) Evaluator {
	return NewEvaluatorWithOptions(scanner, interpreters, Options{NonMutating: true})
}

type evaluator struct {
	scanner      scanner.Scanner
	interpreters Interpreters
	options      Options
}

type Options struct {
	NonMutating   bool
	EvalEmptyTags bool
}

func (ev evaluator) Eval(s, extra interface{}) error {
	return multierror.Prefix(
		ev.eval(
			"",
			nil,
			reflect.ValueOf(s),
			&el.Context{
				Struct:   s,
				Extra:    extra,
				EvalExpr: ev.evalExpr,
			}), "structor:")
}

func (ev evaluator) evalExpr(
	intrprName, expr string,
	ctx *el.Context) (interface{}, error) {
	intrpr, ok := ev.interpreters[intrprName]
	if !ok {
		return nil, fmt.Errorf("unknown interpreter: %s", intrprName)
	}
	return intrpr.Execute(expr, ctx)
}

func (ev evaluator) eval(
	expr string,
	interpreter el.Interpreter,
	v reflect.Value,
	ctx *el.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	if !v.IsValid() {
		return nil
	}
	t := v.Type()
	k := t.Kind()
	elV, elT, elK := v, t, k
	switch k {
	case reflect.Interface, reflect.Ptr:
		elV = v.Elem()
		if elV.IsValid() {
			elT = elV.Type()
			elK = elT.Kind()
		}
	}
	var merr *multierror.Error
	var ctxSub interface{}
	if expr != "" || ev.options.EvalEmptyTags {
		ctx.Val = nil
		if elV.IsValid() {
			ctx.Val = elV.Interface()
		}
		if result, err := interpreter.Execute(expr, ctx); err != nil {
			merr = multierror.Append(merr, err)
		} else {
			ctxSub = result
			if !ev.options.NonMutating {
				nv := result
				if nv == nil {
					if v.IsValid() {
						v.Set(reflect.Zero(t))
					}
				} else {
					vnv := reflect.ValueOf(nv)
					if !vnv.Type().ConvertibleTo(t) &&
						elK == reflect.Struct {
					} else {
						// Try to convert, it may give a panic with suitable
						// message.
						v.Set(vnv.Convert(t))
					}
				}
			}
		}
	}
	switch elK {
	case reflect.Slice, reflect.Array:
		for i, l := 0, elV.Len(); i < l; i++ {
			ctx.LongName = fmt.Sprintf("%s[%d]", ctx.LongName, i)
			ctx.Tags = nil
			v := elV.Index(i)
			err := ev.eval("", nil, v, ctx)
			merr = multierror.Append(merr, err)
		}
	case reflect.Struct:
		ctx.Sub = ctxSub
		for i, l := 0, elV.NumField(); i < l; i++ {
			tf := elT.Field(i)
			tags, err := ev.scanner.Tags(tf.Tag)
			expr := ""
			var interpreter el.Interpreter
			for k, t := range tags {
				if i, ok := ev.interpreters[k]; ok {
					delete(tags, k)
					expr, interpreter = t, i
				}
			}
			if i, ok := ev.interpreters[WholeTag]; ok {
				delete(tags, WholeTag)
				expr, interpreter = string(tf.Tag), i
			}
			ctx.Name = tf.Name
			ctx.LongName = fmt.Sprintf("%s.%s", t, tf.Name)
			ctx.Tags = tags
			v := elV.Field(i)
			err = ev.eval(expr, interpreter, v, ctx)
			merr = multierror.Append(merr, err)
		}
	}
	return multierror.Prefix(
		merr.ErrorOrNil(), fmt.Sprintf("<<%s>>:", ctx.LongName))
}
