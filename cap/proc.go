package cap

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const DefaultMountPoint = "/proc"

type Proc struct {
	PID        int
	MountPoint string
}

// NewProc
func NewProc(pid int) Proc {
	return Proc{
		PID:        pid,
		MountPoint: DefaultMountPoint,
	}
}

func (p Proc) Path(path string) string {
	return filepath.Join(p.MountPoint, strconv.Itoa(p.PID), path)
}

type ProcCaps struct {
	PID    int
	CapInh Caps
	CapPrm Caps
	CapEff Caps
	CapBnd Caps
	CapAmb Caps
}

// NewProcCaps
func NewProcCaps(p Proc) (ProcCaps, error) {
	data, err := ReadFileNoStat(p.Path("status"))
	if err != nil {
		return ProcCaps{}, err
	}
	c := ProcCaps{PID: p.PID}
	clc, err := getCapLastCap(p.MountPoint)
	if err != nil {
		return ProcCaps{}, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if !bytes.Contains([]byte(line), []byte(":")) {
			continue
		}
		kv := strings.SplitN(line, ":", 2)
		k := string(strings.TrimSpace(kv[0]))
		v := string(strings.TrimSpace(kv[1]))
		switch k {
		case "CapInh":
			c.CapInh, err = NewCaps(v, clc)
			if err != nil {
				return ProcCaps{}, err
			}
		case "CapPrm":
			c.CapPrm, err = NewCaps(v, clc)
			if err != nil {
				return ProcCaps{}, err
			}
		case "CapEff":
			c.CapEff, err = NewCaps(v, clc)
			if err != nil {
				return ProcCaps{}, err
			}
		case "CapBnd":
			c.CapBnd, err = NewCaps(v, clc)
			if err != nil {
				return ProcCaps{}, err
			}
		case "CapAmb":
			c.CapAmb, err = NewCaps(v, clc)
			if err != nil {
				return ProcCaps{}, err
			}
		}
	}
	return c, nil
}

// ReadFileNoStat from github.com/prometheus/procfs/internal/util
func ReadFileNoStat(filename string) (b []byte, err error) {
	const maxBufferSize = 1024 * 512

	f, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
	}()

	reader := io.LimitReader(f, maxBufferSize)
	return ioutil.ReadAll(reader)
}
