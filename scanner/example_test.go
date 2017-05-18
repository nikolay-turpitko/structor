package scanner_test

import (
	"fmt"
	"reflect"

	"github.com/nikolay-turpitko/structor/scanner"
)

// ExampleStructTag is an example ported from "reflect" package to
// illustrate support of conventional syntax within tags.
func ExampleStructTag() {
	type S struct {
		F string `species:"gopher" color:"blue"`
	}

	s := S{}
	st := reflect.TypeOf(s)
	field := st.Field(0)
	tags, _ := scanner.Default.Tags(field.Tag)
	fmt.Println(tags["color"], tags["species"])

	// Output:
	// blue gopher
}

// ExampleStructTag_Lookup is an example ported from "reflect" package to
// illustrate support of conventional syntax within tags.
func ExampleStructTag_Lookup() {
	type S struct {
		F0 string `alias:"field_0"`
		F1 string `alias:""`
		F2 string
	}

	s := S{}
	st := reflect.TypeOf(s)
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		tags, _ := scanner.Default.Tags(field.Tag)
		if alias, ok := tags["alias"]; ok {
			if alias == "" {
				fmt.Println("(blank)")
			} else {
				fmt.Println(alias)
			}
		} else {
			fmt.Println("(not specified)")
		}
	}

	// Output:
	// field_0
	// (blank)
	// (not specified)
}

// ExampleStructTag_Ex illustrates use of more relaxed syntax within struct tags.
func ExampleStructTag_Ex() {
	type S struct {
		F string `
# this is an example of multiline tag
# with multiline values, custom quotation
# and field delimiter

species: "gopher"
color:   "blue"

multiline: blue\
gopher\
isn't it strange?

multiline2 = 'yet
another
way'

key-1='value-1', key-2="value-2"`

		F1 string `
			# more convenient formatting possible,
			# when you don't care about spaces in multiline tag values

			species: "gopher" color:"blue"
			key-1: "value-1"; key-2: "value-2"
		`
	}

	s := S{}
	st := reflect.TypeOf(s)
	field := st.Field(0)
	tags, _ := scanner.Default.Tags(field.Tag)
	fmt.Println(tags["color"], tags["species"])
	fmt.Println(tags["multiline"])
	fmt.Println(tags["multiline2"])
	fmt.Println(tags["key-1"], tags["key-2"])

	field = st.Field(1)
	tags, _ = scanner.Default.Tags(field.Tag)
	fmt.Println()
	fmt.Println(tags["color"], tags["species"])
	fmt.Println(tags["key-1"], tags["key-2"])

	// Output:
	// blue gopher
	// blue
	// gopher
	// isn't it strange?
	// yet
	// another
	// way
	// value-1 value-2
	//
	// blue gopher
	// value-1 value-2
}

// ExampleStructTag_Mix illustrate mixed usage of "reflect" package and
// structor.Scanner on the same tags.
//
// Tags for libs, which uses conventional syntax, could be placed on the first
// line and all other lines can be used by the structor.Scanner.
func ExampleStructTag_Mix() {
	type S struct {
		F string `species:"gopher" color:"blue"
				  key-1: "value-1"
				  key-2: "value-2"`
	}

	s := S{}
	st := reflect.TypeOf(s)
	field := st.Field(0)
	fmt.Println(field.Tag.Get("color"), field.Tag.Get("species"))
	// standard lib won't see these fields, but won't panic either
	fmt.Println("'", field.Tag.Get("key-1"), field.Tag.Get("key-2"), "'")
	tags, _ := scanner.Default.Tags(field.Tag)
	fmt.Println(tags["color"], tags["species"])
	fmt.Println(tags["key-1"], tags["key-2"])

	// Output:
	// blue gopher
	// '   '
	// blue gopher
	// value-1 value-2
}
