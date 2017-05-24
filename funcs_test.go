package structor_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nikolay-turpitko/structor"
	"github.com/nikolay-turpitko/structor/el"
	"github.com/nikolay-turpitko/structor/funcs/bytes"
	"github.com/nikolay-turpitko/structor/funcs/crypt"
	"github.com/nikolay-turpitko/structor/funcs/encoding"
	"github.com/nikolay-turpitko/structor/funcs/goquery"
	"github.com/nikolay-turpitko/structor/funcs/math"
	funcs_os "github.com/nikolay-turpitko/structor/funcs/os"
	"github.com/nikolay-turpitko/structor/funcs/regexp"
	"github.com/nikolay-turpitko/structor/funcs/strings"
	"github.com/nikolay-turpitko/structor/funcs/use"
	"github.com/nikolay-turpitko/structor/funcs/xpath"
	"github.com/nikolay-turpitko/structor/scanner"
)

var testEvaluator = structor.NewEvaluator(
	scanner.Default,
	structor.Interpreters{
		structor.WholeTag: &el.DefaultInterpreter{
			AutoEnclose: true,
			Funcs: use.Packages(
				use.Pkg{Prefix: "b_", Funcs: bytes.Pkg},
				use.Pkg{Prefix: "c_", Funcs: crypt.Pkg},
				use.Pkg{Prefix: "e_", Funcs: encoding.Pkg},
				use.Pkg{Prefix: "g_", Funcs: goquery.Pkg},
				use.Pkg{Prefix: "m_", Funcs: math.Pkg},
				use.Pkg{Prefix: "o_", Funcs: funcs_os.Pkg},
				use.Pkg{Prefix: "r_", Funcs: regexp.Pkg},
				use.Pkg{Prefix: "s_", Funcs: strings.Pkg},
				use.Pkg{Prefix: "x_", Funcs: xpath.Pkg},
			),
		},
	})

func TestCrypt(t *testing.T) {
	type theStruct struct {
		A string `c_rot13 "structor"`
		B string `c_rot13 "fgehpgbe"`
		C string `"structor\n" | s_reader | c_md5 | e_hex`
		//D []byte `"KqNVJWMfwUPyVnoPQe5ziXppSa/vIJKcGmbWAHi71LQ=" | e_unbase64 | set`
		D []byte `c_rndKey | set`
		E string `"some plain text message for test"`
		F []byte `.Struct.E | b_bytes | c_aes .Struct.D | set`
		G []byte `.Struct.F | c_unaes .Struct.D | set`
	}
	v := &theStruct{}
	err := testEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, "fgehpgbe", v.A)
	assert.Equal(t, "structor", v.B)
	assert.Equal(t, "a228f448aa0e427fdf2214eb186d2edf", v.C)
	assert.NotNil(t, v.D)
	assert.NotEmpty(t, v.D)
	assert.NotNil(t, v.F)
	assert.NotEmpty(t, v.F)
	assert.NotNil(t, v.G)
	assert.NotEmpty(t, v.G)
	assert.NotEqual(t, v.E, v.F)
	assert.NotEqual(t, v.F, v.G)
	assert.Equal(t, v.E, string(v.G))
}

func TestEncoding(t *testing.T) {
	type theStruct struct {
		A string `e_base64 (b_bytes "structor\n")`
		B []byte `set (e_unbase64 "c3RydWN0b3IK")`
		C string `e_hex (b_bytes "structor")`
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
		A  int `set (m_add 1 2 3 4 5)`
		B1 int `set (m_add (m_mul 2 20) (m_sub 5 3))`
		B2 int `{{with $v := m_sub 5 3}}
				{{m_mul 2 20 | m_add $v | set}}
				{{end}}` // pipes and var
		C int     `set (m_sub 5 2 1)`
		D float64 `m_div 5 2 | set`
	}
	v := &theStruct{}
	err := testEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, 15, v.A)
	assert.Equal(t, 42, v.B1)
	assert.Equal(t, 42, v.B2)
	assert.Equal(t, 2, v.C)
	assert.Equal(t, 2.5, v.D)
}

func TestOS(t *testing.T) {
	type theStruct struct {
		FileNameEnv string `"license_file_name"`

		A string   `o_env .Struct.FileNameEnv`
		B string   `o_readFile "./LICENSE" | set`
		C []byte   `.Struct.FileNameEnv | o_env | o_readFile | set`
		D []string `"./LICENSE" | o_readFile | s_string | s_split "\n" | set`
		E []string `"./LICENSE" | o_readTxtFile | s_split "\n" | set`
	}
	os.Setenv("license_file_name", "./LICENSE")
	v := &theStruct{}
	err := testEvaluator.Eval(v, nil)
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

func TestRegexp(t *testing.T) {
	type theStruct struct {
		A [][]string `"xxx-111-yyy-222" | r_match "(\\w+)-(\\d+)" | set`
		B string     `"because the sky is blue"`
		C string     `index (.Struct.B | r_match "(\\w+)") 2 1`
	}
	v := &theStruct{}
	err := testEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(
		t,
		[][]string{
			{"xxx-111", "xxx", "111"},
			{"yyy-222", "yyy", "222"},
		},
		v.A)
	assert.Equal(t, "sky", v.C)
}

func TestStrings(t *testing.T) {
	type theStruct struct {
		A int      `set (s_atoi "42")`
		B []string `set (s_fields "aaa bbb ccc")`
		C []string `set (s_split "|" "111|222")`
		D string   `s_trimSpace "  xxx  "`
		E string   `"o-!-o" | s_replace "o" "0"`
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

func TestXPath(t *testing.T) {
	extra := struct {
		F1 []byte
		F2 []byte
	}{
		[]byte(`
			<div>
				<span>aaa</span>
				<a href="xxx:bbb">ccc ddd eee fff</a>
				<b class="ggg">   hhh  </b>
			</div>
		`),
		[]byte(`
			<div>
				<span>xxx yyy zzz</span>
				<div class="xxx">
					<a href="http://zzz.com">zzz</a>
					<a href="mailto:zzz@zzz.com">zzz</a>
				</div>
			</div>
		`),
	}
	type theStruct struct {
		// Comment actions before and after text can be used to pass text
		// literally to evaluator with AutoEnclose option.
		// This trick can be useful in case of complex format with quotas etc.
		// Alternatively, data can be passed via .Extra.
		HTML string `{{/* */}}
					<div attr="some attr">zzz</div>
					{{/* */}}`
		A     string `.Extra.F1 | b_reader | x_xpath "//span"`
		B     string `index (.Extra.F1 | b_reader | x_xpath "//a/@href" | s_split ":") 1`
		D     string `index (.Extra.F1 | b_reader | x_xpath "//a" | s_fields) 1`
		E     string `index (.Extra.F1 | b_reader | x_xpath "//a" | r_match "(\\w+)") 2 1`
		H     string `.Extra.F1 | b_reader | x_xpath "//*[@class='ggg']" | s_trimSpace`
		K     string `.Struct.HTML | s_reader | x_xpath "//div/@attr"`
		Email string `index (.Extra.F2 | b_reader | x_xpath "//a[contains(@href,'mailto')]/@href" | s_split ":") 1`
	}
	v := &theStruct{}
	err := testEvaluator.Eval(v, extra)
	assert.NoError(t, err)
	assert.Equal(t, "aaa", v.A)
	assert.Equal(t, "bbb", v.B)
	assert.Equal(t, "ddd", v.D)
	assert.Equal(t, "eee", v.E)
	assert.Equal(t, "hhh", v.H)
	assert.Equal(t, "some attr", v.K)
	assert.Equal(t, "zzz@zzz.com", v.Email)
}

func TestGoquery(t *testing.T) {
	extra := struct {
		F1 []byte
		F2 []byte
	}{
		[]byte(`
			<div>
				<span>aaa</span>
				<a href="xxx:bbb">ccc ddd eee fff</a>
				<b class="ggg">   hhh  </b>
			</div>
		`),
		[]byte(`
			<div>
				<span>xxx yyy zzz</span>
				<div class="xxx">
					<a href="http://zzz.com">zzz</a>
					<a href="mailto:zzz@zzz.com">zzz</a>
				</div>
			</div>
			<h1>header 1</h1>
			<h1>header 2</h1>
			<h1>header 3</h1>
			<h1>header 4</h1>
		`),
	}
	type theStruct struct {
		A      string `(.Extra.F1 | b_reader | g_goquery "span").Text`
		B      string `index ((.Extra.F1 | b_reader | g_goquery "a").First.AttrOr "href" "" | s_split ":") 1`
		H      string `(.Extra.F1 | b_reader | g_goquery ".ggg").First.Text | s_trimSpace`
		Link   string `(.Extra.F2 | b_reader | g_goquery "div.xxx a").First.AttrOr "href" ""`
		Header string `(.Extra.F2 | b_reader | g_goquery "h1:nth-of-type(3)").Text | s_upper`
	}
	v := &theStruct{}
	err := testEvaluator.Eval(v, extra)
	assert.NoError(t, err)
	assert.Equal(t, "aaa", v.A)
	assert.Equal(t, "bbb", v.B)
	assert.Equal(t, "hhh", v.H)
	assert.Equal(t, "http://zzz.com", v.Link)
	assert.Equal(t, "HEADER 3", v.Header)
}

func TestEmbedded(t *testing.T) {
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
			B string `index (index .Sub 0 | s_split ":") 1`
			C string `index (index .Sub 1 | s_split ":") 1`
		} `index (.Extra | s_reader | x_xpath "//span" | r_match "\\[(.+)\\]") 0 1 | s_fields | set`
		B struct {
			D string `.Sub.First.Text`
			E string `.Sub.Last.Text`
		} `.Extra | s_reader | g_goquery "div.xxx p" | set`
		Embedded struct {
			F string `.Extra | s_reader | x_xpath "//b"`
		}
	}
	v := &theStruct{}
	err := testEvaluator.Eval(v, extra)
	assert.NoError(t, err)
	assert.Equal(t, "bbb", v.A.B)
	assert.Equal(t, "ccc", v.A.C)
	assert.Equal(t, "ddd", v.B.D)
	assert.Equal(t, "eee", v.B.E)
	assert.Equal(t, "fff", v.Embedded.F)
}
