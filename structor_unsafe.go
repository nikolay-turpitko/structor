// +build !appengine

package structor

import (
	"reflect"
	"unsafe"
)

func tryUnseal(v reflect.Value) reflect.Value {
	if v.CanAddr() && !v.CanSet() {
		// https://stackoverflow.com/a/43918797/2063744
		return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
	}
	return v
}
