package herd

import (
	"os"
	"testing"
)

func TestDuplicateProviders(t *testing.T) {
	r := NewRegistry("/nonexistent", "/nonexistent")
	p := NewConsulProvider("test")
	p.(*ConsulProvider).Address = "http://consul:8080"
	r.AddProvider(p)
	os.Setenv("CONSUL_HTTP_ADDR", "http://consul:8080")
	r.AddMagicProvider(NewCacheFromProvider(NewConsulProvider("test 2")))
	if len(r.providers) != 1 {
		t.Errorf("AddMagicProviders did not detect duplicate consul provider")
	}
	os.Setenv("CONSUL_HTTP_ADDR", "http://consul:8081")
	r.AddMagicProvider(NewCacheFromProvider(NewConsulProvider("test 3")))
	if len(r.providers) != 2 {
		t.Errorf("AddMagicProviders detected a duplicate provider where there is none")
	}
	os.Unsetenv("CONSUL_HTTP_ADDR")
}
