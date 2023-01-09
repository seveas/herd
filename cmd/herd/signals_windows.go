package main

import (
	"os"

	"github.com/seveas/herd"

	"github.com/sirupsen/logrus"
)

func handleSignals(r *herd.Runner) {
	r.OnSignal(os.Interrupt, func() {
		logrus.Errorf("Interrupted, canceling with unfinished tasks")
		r.Interrupt()
	})
}
