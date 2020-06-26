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
	got, err := c.CapPrm.GetStringSlice()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"CAP_NET_BIND_SERVICE"}
	if diff := cmp.Diff(got, want, nil); diff != "" {
		t.Errorf("got %v\nwant %v", got, want)
	}

	got2, err := c.CapEff.GetStringSlice()
	if err != nil {
		t.Fatal(err)
	}
	want2 := []string{"CAP_NET_BIND_SERVICE"}
	if diff := cmp.Diff(got2, want2, nil); diff != "" {
		t.Errorf("got %v\nwant %v", got2, want2)
	}

	got3, err := c.CapInh.GetStringSlice()
	if err != nil {
		t.Fatal(err)
	}
	want3 := []string{}
	if diff := cmp.Diff(got3, want3, nil); diff != "" {
		t.Errorf("got %v\nwant %v", got3, want3)
	}
}
