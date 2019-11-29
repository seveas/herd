package main

import (
	"os"

	"github.com/seveas/katyusha/cmd/katyusha/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
