package herd

import (
	"fmt"
)

const (
	MAJOR_VERSION = 0
	MINOR_VERSION = 7
	PATCH_VERSION = 1
)

func Version() string {
	return fmt.Sprintf("%d.%d.%d", MAJOR_VERSION, MINOR_VERSION, PATCH_VERSION)
}
