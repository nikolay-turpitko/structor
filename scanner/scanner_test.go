package scanner_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nikolay-turpitko/structor/scanner"
	"github.com/nikolay-turpitko/structor/testhelper"
)

func TestScanner(t *testing.T) {
	fix := []struct {
		fileName string
		expect   map[string]string
	}{
		{
			"./test-1.txt",
			map[string]string{
				"A":        "aaa",
				"B":        "bbb",
				"some-tag": "some-tag",
				"c2_x":     "some data for c2_x",
			},
		},
		{
			"./test-2.txt",
			map[string]string{
				"A":        "aaa",
				"B":        "bbb",
				"some-tag": "some-tag",
				"c2_x":     "some data for c2_x",
			},
		},
		{
			"./test-3.txt",
			map[string]string{
				"A":        "aaa",
				"B":        "bbb",
				"some-tag": "some-tag",
				"c2_x":     "some data for c2_x",
			},
		},
		{
			"./test-4.txt",
			map[string]string{
				"A":        "aaa",
				"B":        "bbb",
				"some-tag": "some-tag",
				"c2_x":     "some data for c2_x",
				"tagX":     "xxx",
				"tagY":     `"yy"`,
				"tagZ":     `\"zz\"`,
			},
		},
		{
			"./test-5.txt",
			map[string]string{
				"A":          "aaa",
				"B":          "bbb",
				"some-tag":   "some-tag",
				"c2_x":       "some data for c2_x",
				"multiline":  "line1 \nline 2 \nline 3",
				"multiline2": "line1 \nline 2 \nline 3",
				"Z":          "some text for Z property",
			},
		},
		{
			"./test-6.txt",
			map[string]string{
				"A":        "aaa",
				"B":        "bbb",
				"some-tag": "some-tag",
				"c2_x":     "some data for c2_x",
				"d-tag":    "line without quotes till endline",
				"aaa":      "alternative delimiters",
				"bbb":      "alternative quotes",
				"$x":       "55",
				"$y":       "77",
			},
		},
	}

	for _, fx := range fix {
		fx := fx
		testhelper.Run(t, fx.fileName, func(t *testing.T) {
			f, err := os.Open(fx.fileName)
			defer f.Close()
			assert.NoError(t, err)
			values, err := scanner.Default.Scan(f)
			assert.NoError(t, err)
			assert.Equal(t, fx.expect, values)
			assert.NotContains(t, values, "")
		})
	}
}
