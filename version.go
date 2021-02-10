package herd

import (
	"fmt"
)

const (
	MAJOR_VERSION = 0
	MINOR_VERSION = 6
	PATCH_VERSION = 0
)

func Version() string {
	return fmt.Sprintf("%d.%d.%d", MAJOR_VERSION, MINOR_VERSION, PATCH_VERSION)
}
