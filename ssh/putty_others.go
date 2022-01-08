//go:build !windows
// +build !windows

package ssh

func (c *config) readPuttyConfig(name string) {
	return
}
