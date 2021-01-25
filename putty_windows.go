package katyusha

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

var PuttyNameMap = map[string]string{}

func puttyConfig(host string) map[string]string {
	ret := make(map[string]string)
	if v, ok := PuttyNameMap[host]; ok {
		host = v
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, fmt.Sprintf(`Software\SimonTatham\PuTTY\Sessions\%s`, host), registry.QUERY_VALUE)
	if err != nil {
		k, err = registry.OpenKey(registry.CURRENT_USER, `Software\SimonTatham\PuTTY\Sessions\Default%20Settings`, registry.QUERY_VALUE)
	}
	if err != nil {
		return ret
	}
	iv, _, err := k.GetIntegerValue("PortNumber")
	if err == nil {
		ret["port"] = fmt.Sprintf("%d", iv)
	}
	sv, _, err := k.GetStringValue("UserName")
	if err == nil {
		ret["user"] = sv
	}
	return ret
}
