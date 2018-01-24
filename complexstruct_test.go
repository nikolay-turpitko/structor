package structor_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nikolay-turpitko/structor"
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
		D: []Inner{Inner{}, Inner{}},
		E: []*Inner{&Inner{}, &Inner{}},
		F: &[]Inner{Inner{}, Inner{}},
		G: &[]*Inner{&Inner{}, &Inner{}},
		H: [2]Inner{Inner{}, Inner{}},
		I: [2]*Inner{&Inner{}, &Inner{}},
		J: &[2]Inner{Inner{}, Inner{}},
		K: &[2]*Inner{&Inner{}, &Inner{}},
		L: [][]Inner{[]Inner{Inner{}, Inner{}}, []Inner{Inner{}, Inner{}}},
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

//TODO: test unaddressable struct
