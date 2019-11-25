package main

import (
	"os"

	"github.com/seveas/katyusha"
	"github.com/seveas/katyusha/cmd/katyusha/cmd"
)

func main() {
	err := cmd.Execute()
	// Make sure all messages have been printed
	if katyusha.UI != nil {
		katyusha.UI.Wait()
	}
	if err != nil {
		os.Exit(1)
	}
}
