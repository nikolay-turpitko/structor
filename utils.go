package structor

import (
	"reflect"

	"github.com/nikolay-turpitko/structor/el"
	"github.com/nikolay-turpitko/structor/scanner"
)

// AddressableCopy returns a pointer to the addressable copy of the struct.
func AddressableCopy(s interface{}) interface{} {
	v := reflect.Indirect(reflect.ValueOf(s))
	s2 := reflect.New(v.Type())
	s2.Elem().Set(v)
	return s2.Interface()
}

// DeepCopy returns a pointer to the deep copy of the struct.
// This is utility function to create a copy of the struct so that
// it would be possible to use the copy with structor without corruption of
// the original struct. It should handle pointers, slices and maps so that
// full independent copy would be created. Copied struct should not have
// any references back to the original struct.
// FIXME: not fully implemented, just a first lame draft.
// TODO: remove this version to separate branch, use some tested and supported
// lib which accurately handles pointers, slices and maps and passes our tests.
func DeepCopy(s interface{}) interface{} {
	c := AddressableCopy(s)
	ev := NewEvaluatorWithOptions(
		scanner.Default,
		Interpreters{
			WholeTag: el.InterpreterFunc(func(s string, ctx *el.Context) (interface{}, error) {
				v := reflect.ValueOf(ctx.Val)
				if !v.IsValid() {
					return ctx.Val, nil
				}
				t, k := v.Type(), v.Kind()
				switch k {
				case reflect.Slice:
					cp := reflect.MakeSlice(t, v.Len(), v.Cap())
					reflect.Copy(cp, v)
					return cp.Interface(), nil
				case reflect.Map:
					cp := reflect.MakeMap(t)
					for _, key := range v.MapKeys() {
						v := reflect.Indirect(v.MapIndex(key))
						if v.IsValid() {
							v2 := reflect.New(v.Type())
							v2.Elem().Set(v)
							cp.SetMapIndex(key, v2)
						}
					}
					return cp.Interface(), nil
				}
				return ctx.Val, nil
			}),
		},
		Options{EvalEmptyTags: true})
	err := ev.Eval(c, nil)
	if err != nil {
		panic(err)
	}
	return c
}
