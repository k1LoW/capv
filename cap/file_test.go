package cap

import (
	"os"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFieCaps(t *testing.T) {
	path := "/tmp/file_cap_test"
	os.Create(path)
	defer os.Remove(path)
	exec.Command("setcap", "cap_net_bind_service=+ep", path).Run()
	c, err := NewFileCaps(path)
	if err != nil {
		t.Fatal(err)
	}
	got := c.CapPrm.GetStringSlice()
	want := []string{"CAP_NET_BIND_SERVICE"}
	if diff := cmp.Diff(got, want, nil); diff != "" {
		t.Errorf("got %v\nwant %v", got, want)
	}

	got2 := c.CapEff
	want2 := FileCapEff(true)
	if got2 != want2 {
		t.Errorf("got %v\nwant %v", got2, want2)
	}

	got3 := c.CapInh.GetStringSlice()
	want3 := []string{}
	if diff := cmp.Diff(got3, want3, nil); diff != "" {
		t.Errorf("got %v\nwant %v", got3, want3)
	}
}
