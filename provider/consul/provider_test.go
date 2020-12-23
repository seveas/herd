package consul

import (
	"testing"
)

func TestProviderEquivalence(t *testing.T) {
	p1 := newConsulProvider("test").(*consulProvider)
	p1.config.Address = "http://consul:8080"

	p2 := newConsulProvider("test 2").(*consulProvider)
	p2.config.Address = "http://consul:8080"
	p2.config.Prefix = "consul:"

	if !p1.Equivalent(p2) {
		t.Errorf("Equivalence not properly detected")
	}

	p2.config.Address = "http://consul:8081"
	if p1.Equivalent(p2) {
		t.Errorf("Non-equivalence not properly detected")
	}
}
