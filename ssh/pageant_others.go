//go:build !windows

package ssh

import "io"

func findPageant() io.ReadWriter {
	return nil
}
