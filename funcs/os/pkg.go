package os

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/nikolay-turpitko/structor/funcs/use"
)

// Pkg contains custom functions defined by this package.
var Pkg = use.FuncMap{
	"env":      os.Getenv,
	"open":     os.Open,
	"readFile": ioutil.ReadFile,
	// func readTxtFile(name string) (string, error)
	// Reads text file into string.
	"readTxtFile": readTxtFile,
	// func exec(name string, arg ...interface{}) ([]byte, error)
	// Executes OS command (process) with given name (path).
	//
	// Convention: if last arg is io.Reader, it goes to stdin of the command.
	"exec": execute,
}

func readTxtFile(name string) (string, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func execute(name string, arg ...interface{}) ([]byte, error) {
	var stdin io.Reader
	hasStdin := false
	args := []string{}
	l := len(arg)
	if l > 0 {
		stdin, hasStdin = arg[len(arg)-1].(io.Reader)
		if hasStdin {
			l--
		}
		for i := 0; i < l; i++ {
			args = append(args, fmt.Sprint(arg[i]))
		}
	}
	cmd := exec.Command(name, args...)
	if hasStdin {
		cmd.Stdin = stdin
	}
	return cmd.Output()
}
