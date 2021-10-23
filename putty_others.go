//go:build !windows
// +build !windows

package herd

func puttyConfig(host string) map[string]string {
	return make(map[string]string)
}
