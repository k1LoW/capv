package cap

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/gommon/color"
)

var Capabilities = [40]string{
	"CAP_CHOWN", // 0
	"CAP_DAC_OVERRIDE",
	"CAP_DAC_READ_SEARCH",
	"CAP_FOWNER",
	"CAP_FSETID",
	"CAP_KILL",
	"CAP_SETGID",
	"CAP_SETUID",
	"CAP_SETPCAP",
	"CAP_LINUX_IMMUTABLE",
	"CAP_NET_BIND_SERVICE",
	"CAP_NET_BROADCAST",
	"CAP_NET_ADMIN",
	"CAP_NET_RAW",
	"CAP_IPC_LOCK",
	"CAP_IPC_OWNER",
	"CAP_SYS_MODULE",
	"CAP_SYS_RAWIO",
	"CAP_SYS_CHROOT",
	"CAP_SYS_PTRACE",
	"CAP_SYS_PACCT",
	"CAP_SYS_ADMIN",
	"CAP_SYS_BOOT",
	"CAP_SYS_NICE",
	"CAP_SYS_RESOURCE",
	"CAP_SYS_TIME",
	"CAP_SYS_TTY_CONFIG",
	"CAP_MKNOD",
	"CAP_LEASE",
	"CAP_AUDIT_WRITE",
	"CAP_AUDIT_CONTROL",
	"CAP_SETFCAP",
	"CAP_MAC_OVERRIDE",
	"CAP_MAC_ADMIN",
	"CAP_SYSLOG",
	"CAP_WAKE_ALARM",
	"CAP_BLOCK_SUSPEND",
	"CAP_AUDIT_READ",
	"CAP_PERFMON",
	"CAP_BPF", // 39
}

type Caps struct {
	caps uint32
	clc  int
}

// NewCaps
func NewCaps(s string, clc int) (Caps, error) {
	caps, err := strconv.ParseUint(string(s), 16, 32)
	if err != nil {
		return Caps{}, err
	}
	return Caps{
		caps: uint32(caps),
		clc:  clc,
	}, nil
}

func (c Caps) Uint32() uint32 {
	return c.caps
}

func (c Caps) GetIntSlice() []int {
	bools := parseCaps(c.caps, c.clc)
	ints := []int{}
	for i, b := range bools {
		if b {
			ints = append(ints, i)
		}
	}
	return ints
}

func (c Caps) GetStringSlice() []string {
	bools := parseCaps(c.caps, c.clc)
	strings := []string{}
	for i, b := range bools {
		if b {
			strings = append(strings, Capabilities[i])
		}
	}
	return strings
}

func (c Caps) String() string {
	return fmt.Sprintf("%s", c.GetStringSlice())
}

type FileCapEff bool

func (c FileCapEff) String() string {
	if c {
		return "1"
	}
	return "0"
}

type NextCaps struct {
	ProcCaps
}

func (c NextCaps) Pretty(w io.Writer) error {
	// P'(permitted)
	if _, err := fmt.Fprintln(w, color.Green("P'(permitted) = (P(inheritable) & F(inheritable)) | (F(permitted) & P(bounding)) | P'(ambient)", color.B)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  %s\n", c.CapPrm.String()); err != nil {
		return err
	}
	// P'(inheritable)
	if _, err := fmt.Fprintln(w, color.Green("P'(inheritable) = P(inheritable)", color.B)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  %s\n", c.CapInh.String()); err != nil {
		return err
	}
	// P'(effective)
	if _, err := fmt.Fprintln(w, color.Green("P'(effective) = F(effective) ? P'(permitted) : P'(ambient)", color.B)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  %s\n", c.CapEff.String()); err != nil {
		return err
	}
	// P'(bounding)
	if _, err := fmt.Fprintln(w, color.Green("P'(bounding) = P(bounding)", color.B)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  %s\n", c.CapBnd.String()); err != nil {
		return err
	}
	// P'(ambient)
	if _, err := fmt.Fprintln(w, color.Green("P'(ambient) = (file is privileged) ? 0 : P(ambient)", color.B)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  %s\n", c.CapAmb.String()); err != nil {
		return err
	}
	return nil
}

func CalcNextCaps(pc ProcCaps, fc FileCaps) NextCaps {
	next := NextCaps{}
	mode := fc.FileInfo.Mode()

	// P'(ambient) = (file is privileged) ? 0 : P(ambient)
	if (mode&os.ModeSetuid != 0) || (mode&os.ModeSetgid != 0) {
		next.CapAmb = Caps{
			caps: 0,
			clc:  pc.CapAmb.clc,
		}
	} else {
		next.CapAmb = pc.CapAmb
	}

	// P'(permitted) = (P(inheritable) & F(inheritable)) |
	//           (F(permitted) & P(bounding)) | P'(ambient)
	next.CapPrm = Caps{
		caps: (pc.CapInh.caps & fc.CapInh.caps) | (fc.CapPrm.caps & pc.CapBnd.caps) | next.CapAmb.caps,
		clc:  pc.CapPrm.clc,
	}

	// P'(effective) = F(effective) ? P'(permitted) : P'(ambient)
	if fc.CapEff {
		next.CapEff = next.CapPrm
	} else {
		next.CapEff = next.CapAmb
	}

	// P'(inheritable) = P(inheritable) [i.e., unchanged]
	next.CapInh = pc.CapInh

	// P'(bounding) = P(bounding) [i.e., unchanged]
	next.CapBnd = pc.CapBnd

	return next
}

func parseCaps(caps uint32, clc int) []bool {
	r := make([]bool, clc+1)
	for i := 0; i <= clc; i++ {
		exp2 := uint32(math.Exp2(float64(i)))
		r[i] = (caps&exp2 > 0)
	}
	return r
}

func getCapLastCap(mt string) (int, error) {
	b, err := ReadFileNoStat(filepath.Join(mt, "sys/kernel/cap_last_cap"))
	if err != nil {
		return -1, err
	}
	i, err := strconv.Atoi(strings.Trim(string(b), "\n"))
	if err != nil {
		return -1, err
	}
	return i, nil
}
