package structor_test

import (
	"os"
	"strings"
	"testing"

	"github.com/apaxa-go/eval"
	"github.com/stretchr/testify/assert"

	"github.com/nikolay-turpitko/structor"
	"github.com/nikolay-turpitko/structor/el/go-el"
	"github.com/nikolay-turpitko/structor/funcs/goquery"
	funcs_os "github.com/nikolay-turpitko/structor/funcs/os"
	"github.com/nikolay-turpitko/structor/funcs/regexp"
	funcs_strings "github.com/nikolay-turpitko/structor/funcs/strings"
	"github.com/nikolay-turpitko/structor/funcs/use"
	"github.com/nikolay-turpitko/structor/funcs/xpath"
	"github.com/nikolay-turpitko/structor/scanner"
)

var testGoEvaluator = structor.NewEvaluator(
	scanner.Default,
	structor.Interpreters{
		structor.WholeTag: &goel.Interpreter{
			Args: eval.ArgsFromInterfaces(use.Packages(
				use.Pkg{Prefix: "goquery.", MapName: strings.Title, Funcs: goquery.Pkg},
				use.Pkg{Prefix: "os.", MapName: strings.Title, Funcs: funcs_os.Pkg},
				use.Pkg{Prefix: "regexp.", MapName: strings.Title, Funcs: regexp.Pkg},
				use.Pkg{Prefix: "strings.", MapName: strings.Title, Funcs: funcs_strings.Pkg},
				use.Pkg{Prefix: "xpath.", MapName: strings.Title, Funcs: xpath.Pkg},
			)),
		},
	})

func TestGoELBasic(t *testing.T) {
	type theStruct struct {
		A int    `1+2+3+4+5`
		B int    `2*20 + (5-3)`
		C string `"4"+"2"`
	}
	v := &theStruct{}
	err := testGoEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, 15, v.A)
	assert.Equal(t, 42, v.B)
	assert.Equal(t, "42", v.C)
}

func TestGoELOS(t *testing.T) {
	type theStruct struct {
		FileNameEnv string `"license_file_name"`

		A string   `os.Env(ctx.Struct.(ctxStruct).FileNameEnv)`
		B string   `os.ReadFile("./LICENSE")`
		C []byte   `os.ReadFile(os.Env(ctx.Struct.(ctxStruct).FileNameEnv))`
		D []string `strings.Split("\n", string(os.ReadFile("./LICENSE")))`
		E []string `strings.Split("\n", os.ReadTxtFile("./LICENSE"))`
	}
	os.Setenv("license_file_name", "./LICENSE")
	v := &theStruct{}
	err := testGoEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, "license_file_name", v.FileNameEnv)
	assert.NotEmpty(t, v.A)
	assert.Equal(t, "./LICENSE", v.A)
	assert.NotEmpty(t, v.B)
	assert.Contains(t, v.B, "MIT")
	assert.Contains(t, string(v.C), "MIT")
	assert.Equal(t, 22, len(v.D))
	assert.Contains(t, v.D[0], "MIT")
	assert.Contains(t, v.D[14], "AS IS")
	assert.Equal(t, 22, len(v.E))
	assert.Contains(t, v.E[0], "MIT")
	assert.Contains(t, v.E[14], "AS IS")
}

func TestGoELRegexp(t *testing.T) {
	type theStruct struct {
		A [][]string `regexp.Match("(\\w+)-(\\d+)", "xxx-111-yyy-222")`
		B string     `"because the sky is blue"`
		C string     `regexp.Match("(\\w+)", ctx.Struct.(ctxStruct).B)[2][1]`
		D string     `regexp.Indx(2, 1, regexp.Match("(\\w+)", ctx.Struct.(ctxStruct).B))`
		E string     `regexp.Indx(0, 1, regexp.Match("(\\d+)", ctx.Struct.(ctxStruct).B))`
	}
	v := &theStruct{}
	err := testGoEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(
		t,
		[][]string{
			{"xxx-111", "xxx", "111"},
			{"yyy-222", "yyy", "222"},
		},
		v.A)
	assert.Equal(t, "sky", v.C)
	assert.Equal(t, "sky", v.D)
	assert.Equal(t, "", v.E)
}

func TestGoELStrings(t *testing.T) {
	type theStruct struct {
		A int      `strings.Atoi("42")`
		B []string `strings.Fields("aaa bbb ccc")`
		C []string `strings.Split("|", "111|222")`
		D string   `strings.TrimSpace("  xxx  ")`
		E string   `strings.Replace("o", "0", "o-!-o")`
		F bool     `strings.Contains("tru", "structor")`
	}
	v := &theStruct{}
	err := testGoEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, 42, v.A)
	assert.Equal(t, []string{"aaa", "bbb", "ccc"}, v.B)
	assert.Equal(t, []string{"111", "222"}, v.C)
	assert.Equal(t, "xxx", v.D)
	assert.Equal(t, "0-!-0", v.E)
	assert.True(t, v.F)
}

func TestGoELEmbedded(t *testing.T) {
	extra := `
		<div>
			<span>some noise [B:bbb C:ccc] more noise</span>
			<div class="xxx">
				<p>ddd</p>
				<p>zzz</p>
				<p>eee</p>
			</div>
			<b>fff</b>
		</div>
		<h1>aaa</h1>
	`
	type theStruct struct {
		A struct {
			B string `strings.Split(":", ctx.Sub.(ctxSub)[0])[1]`
			C string `strings.Split(":", ctx.Sub.(ctxSub)[1])[1]`
		} `strings.Fields(regexp.Match("\\[(.+)\\]", xpath.Xpath("//span", strings.Reader(ctx.Extra.(ctxExtra))))[0][1])`
		B struct {
			D string `ctx.Sub.(ctxSub).First().Text()`
			E string `ctx.Sub.(ctxSub).Last().Text()`
		} `goquery.Goquery("div.xxx p", strings.Reader(ctx.Extra.(ctxExtra)))`
		Embedded struct {
			F string `xpath.Xpath("//b", strings.Reader(ctx.Extra.(ctxExtra)))`
		}
	}
	v := &theStruct{}
	err := testGoEvaluator.Eval(v, extra)
	assert.NoError(t, err)
	assert.Equal(t, "bbb", v.A.B)
	assert.Equal(t, "ccc", v.A.C)
	assert.Equal(t, "ddd", v.B.D)
	assert.Equal(t, "eee", v.B.E)
	assert.Equal(t, "fff", v.Embedded.F)
}
