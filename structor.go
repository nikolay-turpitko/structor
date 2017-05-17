package structor

import (
	"errors"
	"fmt"
	"reflect"
	"text/template"

	multierror "github.com/hashicorp/go-multierror"

	"github.com/nikolay-turpitko/structor/scanner"
)

// ELInterpreter is an interface of EL interpreter.
type ELInterpreter interface {
	Execute(expression string, ctx *Context) (result interface{}, err error)
}

// Context is a context, passed to interpreter.
// It contains information about currently processed field, struct and extra.
type Context struct {
	// Name of the currently processed field.
	Name string
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

// Evaluator is an interface of evaluator, which gets structure and extra
// context as input, iterates over `s`'s fields and evaluate expression tag on
// every field.
type Evaluator interface {
	Eval(s, extra interface{}) error
}

// NewEvaluator returns Evaluator with desired settings.
//
//  evalTag - name of the tag, containing expression for evaluator;
//  interpreter - EL interpreter;
func NewEvaluator(
	evalTag string,
	interpreter ELInterpreter) Evaluator {
	return &evaluator{evalTag, interpreter}
}

// NewDefaultEvaluator returns default Evaluator implementation. Default
// implementation uses tag "eval" for expressions and EL interpreter, based on
// `"text/template"`.
//
//  customFuncs - custom functions, available for interpreter;
func NewDefaultEvaluator(customFuncs template.FuncMap) Evaluator {
	return NewEvaluator(
		"eval", &DefaultInterpreter{AddProvidedFuncs(customFuncs)})
}

type evaluator struct {
	evalTag     string
	interpreter ELInterpreter
}

func (ev evaluator) Eval(s, extra interface{}) error {
	return ev.eval(s, extra, nil, nil)
}

func (ev evaluator) eval(s, extra, substruct, subctx interface{}) error {
	curr := s
	if substruct != nil {
		curr = substruct
	}
	val, typ, err := ev.structIntrospect(curr)
	if err != nil {
		return err
	}
	var merr error
	for i, l := 0, typ.NumField(); i < l; i++ {
		err := func() error {
			expr, name, value, tags, err := ev.fieldIntrospect(val, typ, i)
			if err != nil {
				return fmt.Errorf("<<%T.%s>>: %v", s, name, err)
			}
			if expr == "" {
				return nil
			}
			ctx := &Context{
				Name:   name,
				Val:    value.Interface(),
				Tags:   tags,
				Struct: s,
				Extra:  extra,
				Sub:    subctx,
			}
			result, err := ev.interpreter.Execute(expr, ctx)
			if err != nil {
				return err
			}
			err = reflectSet(value, result)
			if err == nil {
				return nil
			}
			if err != errTryRecursive {
				return fmt.Errorf("<<%T.%s>>: %v", s, name, err)
			}
			return ev.eval(s, extra, value.Interface(), result)
		}()
		if err != nil {
			merr = multierror.Append(merr, err)
		}
	}
	return merr
}

func (ev evaluator) structIntrospect(
	s interface{}) (*reflect.Value, reflect.Type, error) {
	v := reflect.ValueOf(s)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	t := v.Type()
	if t.Kind() != reflect.Struct {
		err := fmt.Errorf(
			"%v must be a struct or a pointer to struct, actually: %v",
			s,
			t.Kind())
		return nil, nil, err
	}
	return &v, t, nil
}

func (ev evaluator) fieldIntrospect(
	val *reflect.Value,
	typ reflect.Type,
	i int) (string, string, *reflect.Value, map[string]string, error) {
	f := typ.Field(i)
	v := val.Field(i)
	tags, err := scanner.Default.Tags(f.Tag)
	if err != nil {
		return "", f.Name, &v, tags, err
	}
	if !v.CanSet() {
		err := fmt.Errorf("%s is not settable", f.Name)
		return "", f.Name, &v, tags, err
	}
	for k, t := range tags {
		if k == ev.evalTag {
			delete(tags, k)
			return t, f.Name, &v, tags, nil
		}
	}
	return "", f.Name, &v, tags, nil
}

var errTryRecursive = errors.New("try recursive") // sentinel error

func reflectSet(pv *reflect.Value, nv interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	vnv := reflect.ValueOf(nv)
	v := *pv
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	vt := v.Type()
	if !vnv.Type().ConvertibleTo(vt) &&
		v.Kind() == reflect.Struct {
		// Try to recursively eval tags on inner struct.
		return errTryRecursive
	}
	// Try to convert, in worst case it'll give a panic with suitable message.
	v.Set(vnv.Convert(vt))
	return nil
}
