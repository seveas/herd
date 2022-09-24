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
	h, _ := p.Load(nil, nil)
	if h[0].Attributes["static_attribute"] != "static_value" {
		t.Errorf("Static attribute has wrong value %s", h[0].Attributes["static_attribute"])
	}
	if h[0].Attributes["dynamic_attribute"] != "dynamic_value_0" {
		t.Errorf("Dynamic attribute has wrong value %s", h[0].Attributes["dynamic_attribute"])
	}
	if h[0].Attributes["config_color"] != "pink" {
		t.Errorf("Configured attribute has wrong value %s", h[0].Attributes["config_color"])
	}
}
