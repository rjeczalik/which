package main

import (
	"fmt"
	"os"

	"github.com/rjeczalik/which"
)

func die(v interface{}) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}

const usage = "usage: gowhich programname"

func main() {
	if len(os.Args) != 2 {
		die(usage)
	}
	c, err := which.Lookup(os.Args[1])
	if err != nil {
		die(err)
	}
	fmt.Println(c.Package)
}
