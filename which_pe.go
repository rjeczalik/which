package which

import (
	"debug/pe"
	"fmt"
)

type filePE struct {
	base uint64
	ptyp *PlatformType
	*pe.File
}

func newpe(path string) (file, error) {
	f, err := pe.Open(path)
	if err != nil {
		return nil, err
	}
	fe := filePE{0, nil, f}
	switch oh := f.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		fe.base = uint64(oh.ImageBase)
		fe.ptyp = PlatformWindows386
	case *pe.OptionalHeader64:
		fe.base = oh.ImageBase
		fe.ptyp = PlatformWindowsAMD64
	default:
		f.Close()
		return nil, ErrNotGoExec
	}
	return fe, nil
}

func (fe filePE) clos()                  { fe.Close() }
func (fe filePE) typ() *PlatformType     { return fe.ptyp }
func (fe filePE) section(string) section { return (section)(nil) }

func newwindowstable(path string, f file) (symtab, pclntab []byte, text uint64, err error) {
	fe := f.(filePE)
	if txt := fe.Section(".text"); txt != nil {
		text = fe.base + uint64(txt.VirtualAddress)
	}
	if pclntab, err = loadPETable(fe.File, "pclntab", "epclntab"); err != nil {
		return
	}
	symtab, err = loadPETable(fe.File, "symtab", "esymtab")
	return
}

// findPESymbol was stolen from $GOROOT/src/cmd/addr2line/main.go:181
func findPESymbol(f *pe.File, name string) (*pe.Symbol, error) {
	for _, s := range f.Symbols {
		if s.Name != name {
			continue
		}
		if s.SectionNumber <= 0 {
			return nil, fmt.Errorf("symbol %s: invalid section number %d", name, s.SectionNumber)
		}
		if len(f.Sections) < int(s.SectionNumber) {
			return nil, fmt.Errorf("symbol %s: section number %d is larger than max %d", name, s.SectionNumber, len(f.Sections))
		}
		return s, nil
	}
	return nil, fmt.Errorf("no %s symbol found", name)
}

// loadPETable was stolen from $GOROOT/src/cmd/addr2line/main.go:197
func loadPETable(f *pe.File, sname, ename string) ([]byte, error) {
	ssym, err := findPESymbol(f, sname)
	if err != nil {
		return nil, err
	}
	esym, err := findPESymbol(f, ename)
	if err != nil {
		return nil, err
	}
	if ssym.SectionNumber != esym.SectionNumber {
		return nil, fmt.Errorf("%s and %s symbols must be in the same section", sname, ename)
	}
	sect := f.Sections[ssym.SectionNumber-1]
	data, err := sect.Data()
	if err != nil {
		return nil, err
	}
	return data[ssym.Value:esym.Value], nil
}
