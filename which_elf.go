package which

import "debug/elf"

type fileElf struct {
	*elf.File
}

type sectionElf struct {
	*elf.Section
}

func (se sectionElf) addr() uint64 {
	return se.Addr
}

func (se sectionElf) data() ([]byte, error) {
	return se.Data()
}

func newelf(path string) (file, error) {
	f, err := elf.Open(path)
	if err != nil {
		return nil, err
	}
	return fileElf{f}, nil
}

func (fe fileElf) clos() {
	fe.Close()
}

func (fe fileElf) typ() (ptyp *PlatformType) {
	switch [2]bool{fe.FileHeader.Class == elf.ELFCLASS64, fe.FileHeader.OSABI == elf.ELFOSABI_FREEBSD} {
	case [2]bool{false, false}:
		ptyp = PlatformLinux386
	case [2]bool{true, false}:
		ptyp = PlatformLinuxAMD64
	case [2]bool{false, true}:
		ptyp = PlatformFreeBSD386
	case [2]bool{true, true}:
		ptyp = PlatformFreeBSDAMD64
	}
	return
}

func (fe fileElf) section(name string) section {
	s := fe.Section("." + name)
	if s == nil {
		return nil
	}
	return sectionElf{s}
}
