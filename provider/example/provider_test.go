package example

import (
	"testing"

	"github.com/spf13/viper"
)

func TestConfigAttribute(t *testing.T) {
	v := viper.New()
	v.SetDefault("Color", "pink")
	p := newProvider("example")
	if err := p.ParseViper(v); err != nil {
		t.Errorf("ParseViper failed: %s", err)
	}
	hs, _ := p.Load(t.Context(), nil)
	h := hs.Get(0)
	if h.Attributes["static_attribute"] != "static_value" {
		t.Errorf("Static attribute has wrong value %s", h.Attributes["static_attribute"])
	}
	if h.Attributes["dynamic_attribute"] != "dynamic_value_0" {
		t.Errorf("Dynamic attribute has wrong value %s", h.Attributes["dynamic_attribute"])
	}
	if h.Attributes["config_color"] != "pink" {
		t.Errorf("Configured attribute has wrong value %s", h.Attributes["config_color"])
	}
}
