//go:build !windows
// +build !windows

package ssh

func (c *configBlock) readPuttyConfig(name string) {
	return
}
