package os

import (
	"io/ioutil"
	"os"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

var Pkg = use.FuncMap{
	"env":      os.Getenv,
	"readFile": ioutil.ReadFile,
	"readTxtFile": func(name string) (string, error) {
		b, err := ioutil.ReadFile(name)
		if err != nil {
			return "", err
		}
		return string(b), nil
	},
}
