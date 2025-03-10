package herd

import (
	"fmt"
)

const (
	majorVersion = 0
	minorVersion = 14
	patchVersion = 1
)

func Version() string {
	return fmt.Sprintf("%d.%d.%d", majorVersion, minorVersion, patchVersion)
}

func VersionTuple() (major, minor, patch int) {
	return majorVersion, minorVersion, patchVersion
}
