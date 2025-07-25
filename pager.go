package herd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type pager struct {
	process *exec.Cmd
	stdin   io.WriteCloser
}

func (p *pager) start() error {
	ctx := context.Background()
	if p == nil || p.process != nil {
		return nil
	}
	pager, ok := os.LookupEnv("PAGER")
	if !ok {
		pager = "less"
	}
	var cmd *exec.Cmd
	if _, err := exec.LookPath(pager); err == nil {
		args := []string{}
		if strings.HasSuffix(pager, "less") {
			args = append(args, "-R")
		}
		cmd = exec.CommandContext(ctx, pager, args...)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", pager)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fd, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	p.process = cmd
	p.stdin = fd
	return cmd.Start()
}

func (p *pager) WriteString(msg string) (int, error) {
	if p.stdin == nil {
		return 0, fmt.Errorf("trying to write to a process that hasn't started")
	}
	return p.stdin.Write([]byte(msg))
}

func (p *pager) Write(msg []byte) (int, error) {
	if p.stdin == nil {
		return 0, fmt.Errorf("trying to write to a process that hasn't started")
	}
	return p.stdin.Write(msg)
}

func (p *pager) Wait() error {
	if p == nil || p.process == nil {
		return nil
	}
	if err := p.stdin.Close(); err != nil {
		return fmt.Errorf("failed to close stdin: %w", err)
	}
	return p.process.Wait()
}
