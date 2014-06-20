package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rjeczalik/which"
)

func die(v interface{}) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}

const usage = "usage: gowhich program_name|executable_path"

func main() {
	if len(os.Args) != 2 {
		die(usage)
	}
	var (
		prog *which.Program
		err  error
	)
	if strings.Contains(os.Args[1], string(os.PathSeparator)) {
		prog, err = which.Look(os.Args[1])
	} else {
		prog, err = which.LookPath(os.Args[1])
	}
	if err != nil {
		die(err)
	}
	fmt.Println(prog.Package)
}
