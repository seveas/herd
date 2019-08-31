package herd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Pager struct {
	Process *exec.Cmd
	Stdin   io.WriteCloser
}

func (p *Pager) Start() error {
	pager, ok := os.LookupEnv("PAGER")
	if !ok {
		pager = "less"
	}
	args := []string{}
	if strings.HasSuffix(pager, "less") {
		args = append(args, "-R")
	}
	cmd := exec.Command(pager, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fd, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	p.Process = cmd
	p.Stdin = fd
	return cmd.Start()
}

func (p *Pager) WriteString(msg string) (int, error) {
	if p.Stdin == nil {
		return 0, fmt.Errorf("trying to write to a process that hasn't started")
	}
	return p.Stdin.Write([]byte(msg))
}

func (p *Pager) Write(msg []byte) (int, error) {
	if p.Stdin == nil {
		return 0, fmt.Errorf("trying to write to a process that hasn't started")
	}
	return p.Stdin.Write(msg)
}

func (p *Pager) Wait() error {
	if p.Process == nil {
		return fmt.Errorf("trying to wait for a process that hasn't started")
	}
	p.Stdin.Close()
	return p.Process.Wait()
}
