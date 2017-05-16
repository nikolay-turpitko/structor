package structor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitTag(t *testing.T) {
	assert.Equal(t, map[string]string{
		"A": "aaa",
		"B": "bbb",
		"C": "ccc",
		"d": "ddd",
		"1": "111",
	}, splitTag(`A:"aaa" B:"bbb" C:"ccc" d:"ddd" 1:"111"`))
}
