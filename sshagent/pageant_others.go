//go:build !windows
// +build !windows

package sshagent

import "io"

func findPageant() io.ReadWriter {
	return nil
}
