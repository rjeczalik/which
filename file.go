package which

import (
	"debug/elf"
	"debug/gosym"
	"errors"
	"runtime"
)

var newfns []func(string) (ExecutableType, []byte, []byte, uint64, error)

func init() {
	switch runtime.GOOS {
	case "darwin":
		newfns = append(newfns, newmacho, newelf, newpe)
	case "windows":
		newfns = append(newfns, newpe, newmacho, newelf)
	default:
		newfns = append(newfns, newelf, newmacho, newpe)
	}
}

// ExecutableType represents the executable file format.
type ExecutableType uint8

const (
	ExecutableElf ExecutableType = iota
	ExecutableMacho
	ExecutablePE
)

// Executable represents a single Go executable file.
type Executable struct {
	Name  string         // Path to the executable.
	Table *gosym.Table   // Go symbol table.
	Type  ExecutableType // Executable file format.
}

// NewExecutable tries to detect executable type for the given path and returns
// a new executable. It fails if file does not exist, is not a Go executable or
// it's unable to parse the file format.
func NewExecutable(name string) (Executable, error) {
	typ, symtab, pclntab, text, err := dumbnew(name)
	if err != nil {
		return Executable{}, err
	}
	lntab := gosym.NewLineTable(pclntab, text)
	if lntab == nil {
		return Executable{}, ErrNotGoExec
	}
	tab, err := gosym.NewTable(symtab, lntab)
	if err != nil {
		return Executable{}, ErrNotGoExec
	}
	return Executable{Name: name, Table: tab, Type: typ}, nil
}

func dumbnew(name string) (typ ExecutableType, symtab, pclntab []byte, text uint64, err error) {
	for _, newfn := range newfns {
		if typ, symtab, pclntab, text, err = newfn(name); err == nil {
			return
		}
	}
	return
}

func newelf(name string) (typ ExecutableType, symtab, pclntab []byte, text uint64, err error) {
	typ = ExecutableElf
	f, err := elf.Open(name)
	if err != nil {
		err = ErrNotGoExec
		return
	}
	defer f.Close()
	sym := f.Section(".gosymtab")
	if sym == nil {
		err = ErrNotGoExec
		return
	}
	symtab, err = sym.Data()
	if err != nil {
		err = ErrNotGoExec
		return
	}
	pcln := f.Section(".gopclntab")
	if pcln == nil {
		err = ErrNotGoExec
		return
	}
	pclntab, err = pcln.Data()
	if err != nil {
		err = ErrNotGoExec
		return
	}
	txt := f.Section(".text")
	if txt != nil {
		err = ErrNotGoExec
		return
	}
	return
}

func newmacho(name string) (typ ExecutableType, symtab, pclntab []byte, text uint64, err error) {
	typ = ExecutableMacho
	// TODO(rjeczalik): osx support
	err = errors.New("not implemented")
	return
}

func newpe(name string) (typ ExecutableType, symtab, pclntab []byte, text uint64, err error) {
	typ = ExecutablePE
	// TODO(rjeczalik): windows support
	err = errors.New("not implemented")
	return
}
