package structor

import (
	"errors"
	"fmt"
	"reflect"
	"text/template"

	multierror "github.com/hashicorp/go-multierror"

	"github.com/nikolay-turpitko/structor/el"
	"github.com/nikolay-turpitko/structor/scanner"
)

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
	interpreter el.Interpreter) Evaluator {
	return &evaluator{evalTag, interpreter}
}

// NewDefaultEvaluator returns default Evaluator implementation. Default
// implementation uses tag "eval" for expressions and EL interpreter, based on
// `"text/template"`.
//
//  customFuncs - custom functions, available for interpreter;
func NewDefaultEvaluator(customFuncs template.FuncMap) Evaluator {
	return NewEvaluator(
		"eval", &el.DefaultInterpreter{AddProvidedFuncs(customFuncs)})
}

type evaluator struct {
	evalTag     string
	interpreter el.Interpreter
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
			longName := fmt.Sprintf("%T.%s", curr, name)
			if err != nil {
				return fmt.Errorf("<<%s>>: %v", longName, err)
			}
			if expr == "" {
				if value.Kind() == reflect.Struct {
					// process embedded struct without tag
					return ev.eval(s, extra, byRef(value), nil)
				}
				return nil
			}
			ctx := &el.Context{
				Name:     name,
				LongName: longName,
				Val:      value.Interface(),
				Tags:     tags,
				Struct:   s,
				Extra:    extra,
				Sub:      subctx,
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
				return fmt.Errorf("<<%s>>: %v", longName, err)
			}
			// process embedded struct with tag
			return ev.eval(s, extra, byRef(value), result)
		}()
		if err != nil {
			merr = multierror.Append(merr, err)
		}
	}
	return merr
}

func (ev evaluator) structIntrospect(
	s interface{}) (reflect.Value, reflect.Type, error) {
	v := indirect(reflect.ValueOf(s))
	t := v.Type()
	if t.Kind() != reflect.Struct {
		err := fmt.Errorf(
			"%v must be a struct or a pointer to struct, actually: %v",
			s,
			t.Kind())
		return v, t, err
	}
	return v, t, nil
}

func (ev evaluator) fieldIntrospect(
	val reflect.Value,
	typ reflect.Type,
	i int) (string, string, reflect.Value, map[string]string, error) {
	f := typ.Field(i)
	v := indirect(val.Field(i))
	tags, err := scanner.Default.Tags(f.Tag)
	if err != nil {
		return "", f.Name, v, tags, err
	}
	if !v.CanSet() && (!v.CanAddr() || !v.Addr().CanSet()) {
		err := fmt.Errorf("%s is not settable", f.Name)
		return "", f.Name, v, tags, err
	}
	for k, t := range tags {
		if k == ev.evalTag {
			delete(tags, k)
			return t, f.Name, v, tags, nil
		}
	}
	return "", f.Name, v, tags, nil
}

var errTryRecursive = errors.New("try recursive") // sentinel error

func reflectSet(v reflect.Value, nv interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	vnv := reflect.ValueOf(nv)
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

func indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

func byRef(v reflect.Value) interface{} {
	if v.CanAddr() {
		v = v.Addr()
	}
	return v.Interface()
}
