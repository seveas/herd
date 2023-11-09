package herd

import (
	"encoding/json"
	"testing"
)

func TestHostDeserialization(t *testing.T) {
	data := map[string]interface{}{
		"Name": "test-host.herd.ci",
		"Attributes": map[string]interface{}{
			"color":  "puce",
			"number": 32,
			"float":  1.1,
		},
	}
	bdata, _ := json.Marshal(data)
	var host Host
	err := json.Unmarshal(bdata, &host)
	if err != nil {
		t.Errorf("Unable to deserialize host data: %s", err)
	}
	if host.Attributes == nil {
		t.Errorf("Deserialized host does not get attributes")
	}
	if domain, ok := host.Attributes["domainname"]; !ok || domain != "herd.ci" {
		t.Errorf("Deserialized host is missing attributes")
	}
	if color, ok := host.Attributes["color"]; !ok || color != "puce" {
		t.Errorf("string attribute did not survive the json trip: %v", host.Attributes)
	}
	if number, ok := host.Attributes["number"]; !ok || number != int64(32) {
		t.Errorf("integer attribute did not survive the json trip: %v", host.Attributes)
	}
	if flt, ok := host.Attributes["float"]; !ok || flt != 1.1 {
		t.Errorf("float attribute did not survive the json trip: %v", host.Attributes)
	}
}

func TestAmendWithoutProvider(t *testing.T) {
	h := NewHost("test-host.herd.ci", "127.0.0.1", HostAttributes{})
	h2 := NewHost("test-host.herd.ci", "127.0.0.1", HostAttributes{})
	h.Amend(h2)
}
