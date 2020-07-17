package main

import (
	"fmt"
	"os"

	"foxygo.at/jsonnext"
	"github.com/google/go-jsonnet"
)

type config struct {
	jsonnext.Config
	Filename string `arg:"" optional:"" help:"File to evaluate. stdin is used if omitted or \"-\""`
}

func main() {
	cli := parseCLI()

	vm := newVM(&cli.Config)

	out, err := run(vm, cli.Filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(out)
}

func newVM(c *jsonnext.Config) *jsonnet.VM {
	vm := jsonnet.MakeVM()
	i := &jsonnext.Importer{}
	c.ConfigureImporter(i, "JXPATH")
	c.ConfigureVM(vm)
	vm.Importer(i)
	return vm
}

func run(vm *jsonnet.VM, filename string) (string, error) {
	if filename == "" || filename == "-" {
		filename = "/dev/stdin"
	}

	node, _, err := vm.ImportAST("", filename)
	if err != nil {
		return "", err
	}

	out, err := vm.Evaluate(node)
	if err != nil {
		return "", err
	}

	return out, nil
}
