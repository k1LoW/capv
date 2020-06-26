package cap

import (
	"math"

	"github.com/syndtr/gocapability/capability"
)

type FileCaps struct {
	Path   string
	CapInh Caps
	CapPrm Caps
	CapEff Caps
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
	return FileCaps{
		Path:   path,
		CapInh: getCaps(c, capability.INHERITABLE, clc),
		CapPrm: getCaps(c, capability.PERMITTED, clc),
		CapEff: getCaps(c, capability.EFFECTIVE, clc),
	}, nil
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
