package structor_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nikolay-turpitko/structor"
	"github.com/nikolay-turpitko/structor/el"
	"github.com/nikolay-turpitko/structor/scanner"
)

func TestArrayAndSlice(t *testing.T) {
	type Inner struct {
		S string `eval:"sss"`
	}
	type T struct {
		A string `eval:"aaa"`
		B Inner
		C *Inner
		D []Inner
		E []*Inner
		F *[]Inner
		G *[]*Inner
		H [2]Inner
		I [2]*Inner
		J *[2]Inner
		K *[2]*Inner
		L [][]Inner
	}
	v := &T{
		C: &Inner{},
		D: []Inner{{}, {}},
		E: []*Inner{{}, {}},
		F: &[]Inner{{}, {}},
		G: &[]*Inner{{}, {}},
		H: [2]Inner{{}, {}},
		I: [2]*Inner{{}, {}},
		J: &[2]Inner{{}, {}},
		K: &[2]*Inner{{}, {}},
		L: [][]Inner{{{}, {}}, {{}, {}}},
	}
	ev := structor.NewDefaultEvaluator(nil)
	err := ev.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, "aaa", v.A)
	assert.Equal(t, "sss", v.B.S)
	assert.Equal(t, "sss", v.C.S)
	assert.Equal(t, "sss", v.D[0].S)
	assert.Equal(t, "sss", v.D[1].S)
	assert.Equal(t, "sss", v.E[0].S)
	assert.Equal(t, "sss", v.E[1].S)
	assert.Equal(t, "sss", (*v.F)[0].S)
	assert.Equal(t, "sss", (*v.F)[1].S)
	assert.Equal(t, "sss", (*v.G)[0].S)
	assert.Equal(t, "sss", (*v.G)[1].S)
	assert.Equal(t, "sss", v.H[0].S)
	assert.Equal(t, "sss", v.H[1].S)
	assert.Equal(t, "sss", v.I[0].S)
	assert.Equal(t, "sss", v.I[1].S)
	assert.Equal(t, "sss", (*v.J)[0].S)
	assert.Equal(t, "sss", (*v.J)[1].S)
	assert.Equal(t, "sss", (*v.K)[0].S)
	assert.Equal(t, "sss", (*v.K)[1].S)
	assert.Equal(t, "sss", v.L[0][0].S)
	assert.Equal(t, "sss", v.L[1][1].S)
}

func TestUnexportedFields(t *testing.T) {
	type Inner struct {
		a string `eval:"aaa"`
	}
	type T struct {
		b string `eval:"bbb"`
		c Inner
	}
	v := &T{c: Inner{}}
	ev := structor.NewDefaultEvaluator(nil)
	err := ev.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, "aaa", v.c.a)
	assert.Equal(t, "bbb", v.b)
}

func TestEmptyTagsAndLongName(t *testing.T) {
	type T struct {
		A string
		S struct {
			B []struct {
				C, d string
			}
		}
		e map[int]*struct{ f string }
	}
	v := &T{
		"",
		struct {
			B []struct{ C, d string }
		}{
			[]struct{ C, d string }{{}, {}},
		},
		map[int]*struct{ f string }{42: {"xxx"}},
	}
	ev := structor.NewEvaluatorWithOptions(
		scanner.Default,
		structor.Interpreters{
			structor.WholeTag: el.InterpreterFunc(func(s string, ctx *el.Context) (interface{}, error) {
				s, ok := ctx.Val.(string)
				if !ok {
					return ctx.Val, nil
				}
				return ctx.LongName, nil
			}),
		},
		structor.Options{EvalEmptyTags: true})
	err := ev.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, "*structor_test.T.A", v.A)
	assert.Equal(t, "*structor_test.T.S.B[0].C", v.S.B[0].C)
	assert.Equal(t, "*structor_test.T.S.B[1].d", v.S.B[1].d)
	assert.Equal(t, "*structor_test.T.e[42].f", v.e[42].f)
}

func TestAddressableCopy(t *testing.T) {
	type Inner struct {
		a string `eval:"aaa"`
	}
	type T struct {
		b string `eval:"bbb"`
		c Inner
	}
	ev := structor.NewDefaultEvaluator(nil)
	v1 := T{
		c: Inner{},
	}
	c1 := structor.AddressableCopy(v1)
	err := ev.Eval(c1, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, "aaa", v1.c.a)
	assert.NotEqual(t, "bbb", v1.b)
	cv1 := c1.(*T)
	assert.Equal(t, "aaa", cv1.c.a)
	assert.Equal(t, "bbb", cv1.b)

	v2 := &T{c: Inner{}}
	c2 := structor.AddressableCopy(v2)
	err = ev.Eval(c2, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, "aaa", v2.c.a)
	assert.NotEqual(t, "bbb", v2.b)
	cv2 := c2.(*T)
	assert.Equal(t, "aaa", cv2.c.a)
	assert.Equal(t, "bbb", cv2.b)
}

func TestDeepCopy(t *testing.T) {
	type Inner struct {
		a string `eval:"aaa"`
	}
	type T struct {
		b string `eval:"bbb"`
		c Inner
		d []Inner
		e map[int]*Inner
	}
	ev := structor.NewDefaultEvaluator(nil)
	v1 := T{
		c: Inner{},
		d: []Inner{{}, {}},
		e: map[int]*Inner{42: &Inner{}},
	}
	c1 := structor.DeepCopy(v1)
	err := ev.Eval(c1, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, "aaa", v1.c.a)
	assert.NotEqual(t, "bbb", v1.b)
	assert.NotEqual(t, "aaa", v1.d[0].a)
	assert.NotEqual(t, "aaa", v1.d[1].a)
	assert.NotEqual(t, "aaa", v1.e[42].a)
	cv1 := c1.(*T)
	assert.Equal(t, "aaa", cv1.c.a)
	assert.Equal(t, "bbb", cv1.b)
	assert.Equal(t, "aaa", cv1.d[0].a)
	assert.Equal(t, "aaa", cv1.d[1].a)
	assert.Equal(t, "aaa", cv1.e[42].a)

	v2 := &T{c: Inner{}}
	c2 := structor.DeepCopy(v2)
	err = ev.Eval(c2, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, "aaa", v2.c.a)
	assert.NotEqual(t, "bbb", v2.b)
	cv2 := c2.(*T)
	assert.Equal(t, "aaa", cv2.c.a)
	assert.Equal(t, "bbb", cv2.b)
}
