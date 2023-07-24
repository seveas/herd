//go:build unix
// +build unix

package main

import (
	"os"
	"syscall"

	"github.com/seveas/herd"

	"github.com/sirupsen/logrus"
)

func handleSignals(r *herd.Runner) {
	r.OnSignal(os.Interrupt, func() {
		logrus.Errorf("Interrupted, canceling with unfinished tasks")
		r.Interrupt()
	})
	r.OnSignal(syscall.SIGUSR1, func() {
		_, s := r.Settings()
		oldp := s["Parallel"].(int)
		p := oldp * 3 / 2
		if p <= 1 || p <= oldp {
			p = oldp + 1
		}
		logrus.Infof("Increasing parallelism to %d", p)
		r.SetParallel(p)
	})
	r.OnSignal(syscall.SIGUSR2, func() {
		_, s := r.Settings()
		p := s["Parallel"].(int) / 2
		if p <= 0 {
			p = 1
		}
		logrus.Infof("Decreasing parallelism to %d", p)
		r.SetParallel(p)
	})
}
