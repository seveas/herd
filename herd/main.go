package main

import (
	"os"

	"github.com/seveas/herd"
	"github.com/seveas/herd/herd/cmd"
)

func main() {
	err := cmd.Execute()
	// Make sure all messages have been printed
	if herd.UI != nil {
		herd.UI.Wait()
	}
	if err != nil {
		os.Exit(1)
	}
}
