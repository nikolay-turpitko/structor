// +build appengine

package structor

import "reflect"

func tryUnseal(v reflect.Value) reflect.Value {
	// unsafe magic is powerless on appengine
	return v
}
