package structor

import (
	"reflect"

	"github.com/mohae/deepcopy"
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
// BUG: Current implementation doesn't copy non-exported fields.
func DeepCopy(s interface{}) interface{} {
	return AddressableCopy(deepcopy.Copy(s))
}
