package cap

import (
	"fmt"
	"io"
	"math"
	"syscall"
	"unsafe"

	"github.com/labstack/gommon/color"
	"github.com/syndtr/gocapability/capability"
)

const vfsCapFlageffective = 0x000001

type FileCaps struct {
	Path   string
	CapInh Caps
	CapPrm Caps
	CapEff FileCapEff
}

func (c FileCaps) Pretty(w io.Writer) error {
	// F(permitted)
	if _, err := fmt.Fprintln(w, color.Yellow("F(permitted):", color.B)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  %s\n", c.CapPrm.String()); err != nil {
		return err
	}
	// F(inheritable)
	if _, err := fmt.Fprintln(w, color.Yellow("F(inheritable):", color.B)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  %s\n", c.CapInh.String()); err != nil {
		return err
	}
	// F(effective)
	if _, err := fmt.Fprintln(w, color.Yellow("F(effective):", color.B)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  %s\n", c.CapEff.String()); err != nil {
		return err
	}
	return nil
}

// NewFileCaps
func NewFileCaps(path string) (FileCaps, error) {
	c, err := capability.NewFile2(path)
	if err != nil {
		return FileCaps{}, err
	}
	err = c.Load()
	if err != nil {
		return FileCaps{}, err
	}
	clc, err := getCapLastCap(DefaultMountPoint)
	if err != nil {
		return FileCaps{}, err
	}
	capEff, err := getFileCapEffective(path)
	if err != nil {
		return FileCaps{}, err
	}
	return FileCaps{
		Path:   path,
		CapInh: getCaps(c, capability.INHERITABLE, clc),
		CapPrm: getCaps(c, capability.PERMITTED, clc),
		CapEff: capEff,
	}, nil
}

func getFileCapEffective(path string) (FileCapEff, error) {
	name := "security.capability"
	var v uint32
	_, _, e1 := syscall.Syscall6(syscall.SYS_GETXATTR, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))), uintptr(unsafe.Pointer(syscall.StringBytePtr(name))), uintptr(unsafe.Pointer(&v)), uintptr(4*(1+2*2)), 0, 0)
	if e1 != 0 {
		if e1 == syscall.ENODATA {
			return FileCapEff(false), nil
		}
		return FileCapEff(false), e1
	}
	return FileCapEff(v&vfsCapFlageffective != 0), nil
}

// getCaps
func getCaps(c capability.Capabilities, which capability.CapType, clc int) Caps {
	var caps uint32
	for i := 0; i <= clc; i++ {
		var cap uint32
		if c.Get(which, capability.Cap(i)) {
			cap = 1
		}
		caps = caps + cap*uint32(math.Pow(2, float64(i)))
	}
	return Caps{
		caps: caps,
		clc:  clc,
	}
}
