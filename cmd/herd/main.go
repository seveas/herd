package main

import (
	"os"

	"github.com/seveas/herd/cmd/herd/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
