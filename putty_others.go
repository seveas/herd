// +build !windows

package herd

import (
	"io"
)

func findPageant() io.ReadWriter {
	return nil
}

func puttyConfig(host string) map[string]string {
	return make(map[string]string)
}
