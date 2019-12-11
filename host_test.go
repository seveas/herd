package herd

import (
	"encoding/json"
	"testing"
)

func TestHostDeserialization(t *testing.T) {
	data := map[string]interface{}{
		"Name": "test-host",
	}
	bdata, _ := json.Marshal(data)
	var host Host
	err := json.Unmarshal(bdata, &host)
	if err != nil {
		t.Errorf("Unable to deserialize host data: %s", err)
	}
	host.init()
	if host.Attributes == nil {
		t.Errorf("Deserialized host does not get attributes")
	}
	if host.Port != 22 {
		t.Errorf("Deserialized host does not get default port")
	}
}
