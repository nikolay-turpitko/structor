package structor_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nikolay-turpitko/structor"
	"github.com/nikolay-turpitko/structor/el"
	"github.com/nikolay-turpitko/structor/funcs/crypt"
	"github.com/nikolay-turpitko/structor/funcs/encoding"
	"github.com/nikolay-turpitko/structor/funcs/math"
	"github.com/nikolay-turpitko/structor/funcs/os"
	"github.com/nikolay-turpitko/structor/funcs/regexp"
	"github.com/nikolay-turpitko/structor/funcs/strings"
	"github.com/nikolay-turpitko/structor/funcs/use"
)

var testEvaluator = structor.NewEvaluator(structor.Interpreters{
	structor.WholeTag: &el.DefaultInterpreter{
		AutoEnclose: true,
		Funcs: use.Packages(
			use.Pkg{"c_", crypt.Pkg},
			use.Pkg{"e_", encoding.Pkg},
			use.Pkg{"m_", math.Pkg},
			use.Pkg{"o_", os.Pkg},
			use.Pkg{"r_", regexp.Pkg},
			use.Pkg{"s_", strings.Pkg},
		),
	},
})

func TestCrypt(t *testing.T) {
	type theStruct struct {
		A string `c_rot13 "structor"`
		B string `c_rot13 "fgehpgbe"`
	}
	v := &theStruct{}
	err := testEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, "fgehpgbe", v.A)
	assert.Equal(t, "structor", v.B)
}

func TestEncoding(t *testing.T) {
	type theStruct struct {
		A string `e_base64 (e_bytes "structor\n")`
		B []byte `set (e_unbase64 "c3RydWN0b3IK")`
		C string `e_hex (e_bytes "structor")`
		D []byte `set (e_unhex "7374727563746f72")`
	}
	v := &theStruct{}
	err := testEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, "c3RydWN0b3IK", v.A)
	assert.Equal(t, []byte("structor\n"), v.B)
	assert.Equal(t, "7374727563746f72", v.C)
	assert.Equal(t, []byte("structor"), v.D)
}

func TestMath(t *testing.T) {
	type theStruct struct {
		A int     `set (m_add 1 2 3 4 5)`
		B int     `set (m_add (m_mul 2 20)  (m_sub 5 3))`
		C int     `set (m_sub 5 2 1)`
		D float64 `set (m_div 5 2)`
	}
	v := &theStruct{}
	err := testEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, 15, v.A)
	assert.Equal(t, 42, v.B)
	assert.Equal(t, 2, v.C)
	assert.Equal(t, 2.5, v.D)
}

func TestOS(t *testing.T) {
	type theStruct struct {
		A string   `o_env "GOROOT"`
		B string   `set (o_readFile "./LICENSE")`
		C []byte   `set (o_readFile "./LICENSE")`
		D []string `set (s_split (s_string (o_readFile "./LICENSE")) "\n")`
		E []string `set (s_split (o_readTxtFile "./LICENSE") "\n")`
	}
	v := &theStruct{}
	err := testEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, v.A)
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

func TestRegexp(t *testing.T) {
	type theStruct struct {
		A [][]string `set (r_match "xxx-111-yyy-222" "(\\w+)-(\\d+)")`
	}
	v := &theStruct{}
	err := testEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, [][]string{{"xxx-111", "xxx", "111"}, {"yyy-222", "yyy", "222"}}, v.A)
}

func TestStrings(t *testing.T) {
	type theStruct struct {
		A int      `set (s_atoi "42")`
		B []string `set (s_fields "aaa bbb ccc")`
		C []string `set (s_split "111|222" "|")`
		D string   `s_trimSpace "  xxx  "`
		E string   `s_replace "o-!-o" "o" "0"`
	}
	v := &theStruct{}
	err := testEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, 42, v.A)
	assert.Equal(t, []string{"aaa", "bbb", "ccc"}, v.B)
	assert.Equal(t, []string{"111", "222"}, v.C)
	assert.Equal(t, "xxx", v.D)
	assert.Equal(t, "0-!-0", v.E)
}
