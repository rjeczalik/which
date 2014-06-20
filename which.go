package which

import (
	"debug/elf"
	"debug/gosym"
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
)

type Cmd struct {
	Path    string
	Package string
}

var errNotFound = errors.New("error whichgo: not found")
var errNotGoExec = errors.New("error whichgo: not a Go executable")
var errGuessFail = errors.New("error whichgo: unable to guess package main import path")

func Lookup(name string) (*Cmd, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil, errNotFound
	}
	f, err := elf.Open(path)
	if err != nil {
		return nil, errNotGoExec
	}
	defer f.Close()
	sym := f.Section(".gosymtab")
	if sym == nil {
		return nil, errNotGoExec
	}
	symdat, err := sym.Data()
	if err != nil {
		return nil, errNotGoExec
	}
	pcln := f.Section(".gopclntab")
	if pcln == nil {
		return nil, errNotGoExec
	}
	pclndat, err := pcln.Data()
	if err != nil {
		return nil, errNotGoExec
	}
	text := f.Section(".text")
	if text == nil {
		return nil, errNotGoExec
	}
	lntab := gosym.NewLineTable(pclndat, text.Addr)
	if lntab == nil {
		return nil, errNotGoExec
	}
	tab, err := gosym.NewTable(symdat, lntab)
	if err != nil {
		return nil, errNotGoExec
	}
	var pkg string
	for file := range tab.Files {
		if strings.Contains(file, "cmd/"+name) {
			if i := strings.LastIndex(file, "/src/"); i != -1 {
				pkg = filepath.Dir(file)[i+len("/src/"):]
				break
			}
		}
	}
	if pkg == "" {
		return nil, errGuessFail
	}
	return &Cmd{Path: path, Package: pkg}, nil
}
