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
	"unsafe"

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

// NewEvaluatorWithOptions returns Evaluator with desired settings.
//
// Only first tag with EL will be recognized and used (only one
// expression per struct field). Different fields of the same struct can be
// processed using different EL interpreters.
//
//  scanner - is a scanner implementation to be used to scan tags.
//  interpreters - is a map of registered tag names to EL interpreters.
//  options - is an Options structure.
func NewEvaluatorWithOptions(
	scanner scanner.Scanner,
	interpreters Interpreters,
	options Options) Evaluator {
	if len(interpreters) == 0 {
		panic("no interpreters registered")
	}
	return &evaluator{scanner, interpreters, options}
}

// NewEvaluator returns Evaluator with desired settings.
// It invokes NewEvaluatorWithOptions with default options.
func NewEvaluator(
	scanner scanner.Scanner,
	interpreters Interpreters) Evaluator {
	return NewEvaluatorWithOptions(scanner, interpreters, Options{})
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
// change original structure (does not save evaluated results) itself.  Though,
// interpreters can change structures' fields as a side effect.  In this case
// Evaluator can be used as a visitor of fields with tags for which it have
// registered interpreters.  It will invoke registered interpreter for field
// with corresponded tag.  Interpreter can then manipulate it's own state or
// el.Context.  For example, it can store processing results into context's
// Extra field.
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

// Options is an options to create Evaluator.
type Options struct {

	// NonMutating creates non-mutating Evaluator, see NewNonmutatingEvaluator.
	NonMutating bool

	// EvalEmptyTags causes Evaluator to invoke Interpreter for fields with
	// empty tags.
	EvalEmptyTags bool
}

func (ev evaluator) Eval(s, extra interface{}) error {
	v := reflect.ValueOf(s)
	t := v.Type()
	k := t.Kind()
	prev := v
	for v.IsValid() && (k == reflect.Interface || k == reflect.Ptr) {
		prev = v
		v = v.Elem()
		if v.IsValid() {
			t = v.Type()
			k = t.Kind()
		}
		if prev == v {
			break
		}
	}
	if k != reflect.Struct || !v.CanSet() {
		return fmt.Errorf("structor: %T: not a settable struct", s)
	}
	return multierror.Prefix(
		ev.eval(
			"",
			nil,
			reflect.ValueOf(s),
			&el.Context{
				Struct:   s,
				Extra:    extra,
				EvalExpr: ev.evalExpr,
				LongName: fmt.Sprintf("%T", s),
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
	// Note: only errors returned by recursive call should not be prefixed.
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = multierror.Prefix(
					fmt.Errorf("%v", r),
					fmt.Sprintf("<<%s>>", ctx.LongName))
			}
		}
	}()
	if !v.IsValid() {
		return nil
	}
	t := v.Type()
	k := t.Kind()
	if v.CanAddr() && !v.CanSet() {
		// https://stackoverflow.com/a/43918797/2063744
		v = reflect.NewAt(t, unsafe.Pointer(v.UnsafeAddr())).Elem()
	}
	elV, elT, elK := v, t, k
	for elV.IsValid() && (elK == reflect.Interface || elK == reflect.Ptr) {
		v, t, k = elV, elT, elK
		elV = v.Elem()
		if elV.IsValid() {
			elT = elV.Type()
			elK = elT.Kind()
		}
		if v == elV {
			break
		}
	}
	if elV.IsValid() {
		if elV.CanAddr() && !elV.CanSet() {
			// https://stackoverflow.com/a/43918797/2063744
			elV = reflect.NewAt(elT, unsafe.Pointer(elV.UnsafeAddr())).Elem()
		}
		elT = elV.Type()
		elK = elT.Kind()
	}
	var merr *multierror.Error
	var ctxSub interface{}
	if (expr != "" || ev.options.EvalEmptyTags) && interpreter != nil {
		ctx.Val = nil
		if elV.IsValid() {
			ctx.Val = elV.Interface()
		}
		if result, err := interpreter.Execute(expr, ctx); err != nil {
			merr = multierror.Append(
				merr,
				multierror.Prefix(err, fmt.Sprintf("<<%s>>", ctx.LongName)))
		} else {
			ctxSub = result
			if !ev.options.NonMutating {
				if result == nil {
					if v.IsValid() {
						v.Set(reflect.Zero(t))
					}
				} else {
					vnv := reflect.ValueOf(result)
					if vnv.Type().ConvertibleTo(t) || elK != reflect.Struct {
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
		prevLongName := ctx.LongName
		for i, l := 0, elV.Len(); i < l; i++ {
			v := elV.Index(i)
			ctx.Name = v.Type().Name()
			ctx.LongName = fmt.Sprintf("%s[%d]", ctx.LongName, i)
			ctx.Tags = nil
			err := ev.eval("", nil, v, ctx)
			ctx.LongName = prevLongName
			merr = multierror.Append(merr, err)
		}
	case reflect.Struct:
		ctx.Sub = ctxSub
		prevLongName := ctx.LongName
		for i, l := 0, elV.NumField(); i < l; i++ {
			tf := elT.Field(i)
			tags, err := ev.scanner.Tags(tf.Tag)
			if err != nil {
				merr = multierror.Append(
					merr,
					multierror.Prefix(err, fmt.Sprintf("<<%s>>", ctx.LongName)))
				break
			}
			expr := ""
			var interpreter el.Interpreter
			for k, t := range tags {
				if i, ok := ev.interpreters[k]; ok {
					delete(tags, k)
					expr, interpreter = t, i
				}
			}
			if i, ok := ev.interpreters[WholeTag]; ok && interpreter == nil {
				delete(tags, WholeTag)
				expr, interpreter = string(tf.Tag), i
			}
			ctx.Name = tf.Name
			ctx.LongName = fmt.Sprintf("%s.%s", ctx.LongName, tf.Name)
			ctx.Tags = tags
			v := elV.Field(i)
			err = ev.eval(expr, interpreter, v, ctx)
			ctx.LongName = prevLongName
			merr = multierror.Append(merr, err)
		}
	}
	return merr.ErrorOrNil()
}

// AddressableCopy returns a pointer to the addressable copy of the struct.
func AddressableCopy(s interface{}) interface{} {
	v := reflect.Indirect(reflect.ValueOf(s))
	s2 := reflect.New(v.Type())
	s2.Elem().Set(v)
	return s2.Interface()
}
