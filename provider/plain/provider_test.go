package plain

import (
	"testing"

	"github.com/seveas/herd"
)

func TestRelativeFiles(t *testing.T) {
	r := herd.NewRegistry("/foo", "/bar")
	p := &plainTextProvider{name: "ittest"}
	p.config.File = "inventory"
	r.AddProvider(p)
	if p.config.File != "/foo/inventory" {
		t.Errorf("Filepath did not get interpreted relative to dataDir")
	}
}
