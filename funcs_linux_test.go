package structor_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	type theStruct struct {
		A string `"./LICENSE"`

		// Dynamically evaluate shell script.
		B int `printf "cat %s | wc -l" .Struct.A | o_exec "/bin/sh" "-c" | s_string | s_trimSpace | s_atoi | set`

		// Pipe from EL expression to shell script.
		C int `o_open .Struct.A | o_exec "/bin/sh" "-c" "wc -l" | s_string | s_trimSpace | s_atoi | set`

		// Test with a reader. Note that reader can be used only once.
		D io.Reader `o_open .Struct.A | set`
		E int       `.Struct.D | o_exec "/bin/sh" "-c" "wc -l" | s_string | s_trimSpace | s_atoi | set`
		F int       `.Struct.D | o_readAll | s_string | len | set`
		G io.Reader `o_open .Struct.A | set`
		H int       `.Struct.G | o_readAll | s_string | len | set`
	}
	v := &theStruct{}
	err := testEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, "./LICENSE", v.A)
	assert.Equal(t, 21, v.B)
	assert.Equal(t, 21, v.C)
	assert.Equal(t, 21, v.E)
	assert.Equal(t, 0, v.F)
	assert.Equal(t, 1073, v.H)
}
