package herd

import (
	"testing"
)

func TestNilSshConfig(t *testing.T) {
	var c *sshConfig
	hc := c.configForHost("test-host")
	if len(hc) != 0 {
		t.Errorf("Expected an empty dict from <nil>.configForhost, got %v", hc)
	}
}
