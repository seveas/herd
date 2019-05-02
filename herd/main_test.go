package main

import (
	"os"
	"testing"
)

func TestMainfunc(t *testing.T) {
	os.Args = []string{"herd"}
	main()
}
