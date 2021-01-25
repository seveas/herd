// +build !windows

package katyusha

func puttyConfig(host string) map[string]string {
	return make(map[string]string)
}
