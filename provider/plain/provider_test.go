package plain

import (
	"testing"

	"github.com/seveas/katyusha"
)

func TestRelativeFiles(t *testing.T) {
	r := katyusha.NewRegistry("/foo", "/bar")
	p := &plainTextProvider{name: "ittest"}
	p.config.File = "inventory"
	r.AddProvider(p)
	if p.config.File != "/foo/inventory" {
		t.Errorf("Filepath did not get interpreted relative to dataDir")
	}
}
