package cap

import (
	"math"
	"path/filepath"
	"strconv"
	"strings"
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

func (c Caps) GetIntSlice() ([]int, error) {
	bools := parseCaps(c.caps, c.clc)
	ints := []int{}
	for i, b := range bools {
		if b {
			ints = append(ints, i)
		}
	}
	return ints, nil
}

func (c Caps) GetStringSlice() ([]string, error) {
	bools := parseCaps(c.caps, c.clc)
	strings := []string{}
	for i, b := range bools {
		if b {
			strings = append(strings, Capabilities[i])
		}
	}
	return strings, nil
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
