// +build go1.7

package structor_test

import (
	"fmt"
	"strings"

	"github.com/apaxa-go/eval"

	"github.com/nikolay-turpitko/structor"
	"github.com/nikolay-turpitko/structor/el"
	goel "github.com/nikolay-turpitko/structor/el/go-el"
	funcs_strings "github.com/nikolay-turpitko/structor/funcs/strings"
	"github.com/nikolay-turpitko/structor/funcs/use"
	"github.com/nikolay-turpitko/structor/scanner"
)

// Example_complex is a more complex example of structor's usage.
//
// Two interpreters are registered to process correspondent struct's tags.
// Custom strings functions are registered with custom prefix and name translations.
// EL uses value from other tag via execution context.
// One interpreter calls another to evaluate expression, got from some tag.
// Custom quotation marks are used in tags, when it's convenient.
func Example_complex() {
	ev := structor.NewEvaluator(
		scanner.Default,
		structor.Interpreters{
			"tmplEL": &el.DefaultInterpreter{
				AutoEnclose: true,
				Funcs:       use.Packages(use.Pkg{Prefix: "str", MapName: strings.Title, Funcs: funcs_strings.Pkg}),
			},
			"goEL": &goel.Interpreter{
				Args: eval.ArgsFromInterfaces(use.Packages(
					use.Pkg{Prefix: "strings.", MapName: strings.Title, Funcs: funcs_strings.Pkg},
				)),
			},
		})
	type theStruct struct {
		A string `tmplEL:".Tags.arg | strUpper" arg:"structor"`
		B string `goEL:'strings.Upper(ctx.Tags["arg"])' arg:"structor"`
		C int    `tmplEL:'.Tags.expr | eval "goEL" | set' expr:"40+2"`
		D int    `goEL:'eval("tmplEL", ctx.Tags["expr"])' expr:'"42" | strAtoi | set'`
	}
	v := &theStruct{}
	err := ev.Eval(v, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(v.A)
	fmt.Println(v.B)
	fmt.Println(v.C)
	fmt.Println(v.D)

	// Output: STRUCTOR
	// STRUCTOR
	// 42
	// 42
}
