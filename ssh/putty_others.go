//go:build !windows
// +build !windows

package ssh

func (c *Config) readPuttyConfig(name string) {
	return
}
