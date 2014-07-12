package which

import "debug/macho"

type fileMacho struct {
	*macho.File
}

type sectionMacho struct {
	*macho.Section
}

func (sm sectionMacho) addr() uint64 {
	return sm.Addr
}

func (sm sectionMacho) data() ([]byte, error) {
	return sm.Data()
}

func newmacho(path string) (file, error) {
	f, err := macho.Open(path)
	if err != nil {
		return nil, err
	}
	fe := fileMacho{f}
	return fe, nil
}

func (fm fileMacho) clos() {
	fm.Close()
}

func (fm fileMacho) typ() (ptyp *PlatformType) {
	switch fm.Cpu {
	case macho.Cpu386:
		ptyp = PlatformDarwin386
	case macho.CpuAmd64:
		ptyp = PlatformDarwinAMD64
	}
	return
}

func (fm fileMacho) section(name string) section {
	s := fm.Section("__" + name)
	if s == nil {
		return nil
	}
	return sectionMacho{s}
}
