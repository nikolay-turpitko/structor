package structor_test

import (
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
	}
	v := &theStruct{}
	err := testEvaluator.Eval(v, nil)
	assert.NoError(t, err)
	assert.Equal(t, "./LICENSE", v.A)
	assert.Equal(t, 21, v.B)
}
