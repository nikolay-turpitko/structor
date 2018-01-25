// +build !appengine

package structor_test

import (
	"testing"

	"github.com/nikolay-turpitko/structor"
	"github.com/stretchr/testify/assert"
)

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
		A string `eval:"aaa"`
	}
	type T struct {
		B string `eval:"bbb"`
		C Inner
		D []Inner
		E map[int]*Inner
	}
	ev := structor.NewDefaultEvaluator(nil)
	v1 := T{
		C: Inner{},
		D: []Inner{{}, {}},
		E: map[int]*Inner{42: &Inner{}},
	}
	c1 := structor.DeepCopy(v1)
	err := ev.Eval(c1, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, "aaa", v1.C.A)
	assert.NotEqual(t, "bbb", v1.B)
	assert.NotEqual(t, "aaa", v1.D[0].A)
	assert.NotEqual(t, "aaa", v1.D[1].A)
	assert.NotEqual(t, "aaa", v1.E[42].A)
	cv1 := c1.(*T)
	assert.Equal(t, "aaa", cv1.C.A)
	assert.Equal(t, "bbb", cv1.B)
	assert.Equal(t, "aaa", cv1.D[0].A)
	assert.Equal(t, "aaa", cv1.D[1].A)
	assert.Equal(t, "aaa", cv1.E[42].A)

	v2 := &T{C: Inner{}}
	c2 := structor.DeepCopy(v2)
	err = ev.Eval(c2, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, "aaa", v2.C.A)
	assert.NotEqual(t, "bbb", v2.B)
	cv2 := c2.(*T)
	assert.Equal(t, "aaa", cv2.C.A)
	assert.Equal(t, "bbb", cv2.B)
}
