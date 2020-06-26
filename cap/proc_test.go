package cap

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const testMountPoint = "../testdata/proc"

func TestProcCaps(t *testing.T) {
	p := Proc{
		PID:        1,
		MountPoint: testMountPoint,
	}

	c, err := NewProcCaps(p)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.CapPrm.GetStringSlice()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"CAP_CHOWN", "CAP_DAC_OVERRIDE", "CAP_FOWNER", "CAP_FSETID", "CAP_KILL", "CAP_SETGID", "CAP_SETUID", "CAP_SETPCAP", "CAP_NET_BIND_SERVICE", "CAP_NET_RAW", "CAP_SYS_CHROOT", "CAP_MKNOD", "CAP_AUDIT_WRITE", "CAP_SETFCAP"}

	if diff := cmp.Diff(got, want, nil); diff != "" {
		t.Errorf("got %v\nwant %v", got, want)
	}

	got2, err := c.CapPrm.GetIntSlice()
	if err != nil {
		t.Fatal(err)
	}
	want2 := []int{0, 1, 3, 4, 5, 6, 7, 8, 10, 13, 18, 27, 29, 31}

	if diff := cmp.Diff(got2, want2, nil); diff != "" {
		t.Errorf("got %v\nwant %v", got2, want2)
	}
}
