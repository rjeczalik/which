package which

import "testing"

const echo = "github.com/rjeczalik/which/testdata/cmd/echo"

var testdata = map[*PlatformType]string{
	PlatformDarwin386:    "testdata/darwin_386/echo",
	PlatformDarwinAMD64:  "testdata/darwin_amd64/echo",
	PlatformFreeBSD386:   "testdata/freebsd_386/echo",
	PlatformFreeBSDAMD64: "testdata/freebsd_amd64/echo",
	PlatformLinux386:     "testdata/linux_386/echo",
	PlatformLinuxAMD64:   "testdata/linux_amd64/echo",
	// TODO(rjeczalik): #2
	// PlatformWindows386:   "testdata/windows_386/echo.exe",
	// PlatformWindowsAMD64: "testdata/windows_amd64/echo.exe",
}

func TestNewExec(t *testing.T) {
	for etyp, path := range testdata {
		ex, err := NewExec(path)
		if err != nil {
			t.Errorf("want err=nil; got %q (etyp=%v)", err, etyp)
			continue
		}
		if ex.Type != etyp {
			t.Errorf("want ex.Type=%v; got %v", etyp, ex.Type)
		}
	}
}

func TestImport(t *testing.T) {
	for etyp, path := range testdata {
		ex, err := NewExec(path)
		if err != nil {
			t.Errorf("want err=nil; got %q (etyp=%v)", err, etyp)
			continue
		}
		imp, err := ex.Import()
		if err != nil {
			t.Errorf("want err=nil; got %q (etyp=%v)", err, etyp)
			continue
		}
		if imp != echo {
			t.Errorf("want imp=%q; got %q (etyp=%v)", echo, imp, etyp)
		}
	}
}
