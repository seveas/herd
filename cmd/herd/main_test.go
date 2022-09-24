package main

import (
	"os"
	"testing"
)

func TestMainfunc(t *testing.T) {
	_ = t
	os.Args = []string{"herd"}
	main()
}
