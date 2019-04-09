package herd

import (
	"fmt"
	"time"
)

type Result struct {
	Host   string
	Stdout []byte
	Stderr []byte
	Err    error
}

func (r Result) String() string {
	return fmt.Sprintf("[%s] (Err: %s)]\n%s\n---\n%s\n", r.Host, r.Err, string(r.Stdout), string(r.Stderr))
}

type HistoryItem struct {
	Hosts   []*Host
	Command string
	Results map[string]Result
}

type History []HistoryItem

type RunnerConfig struct {
	Timeout time.Duration
}

type Runner struct {
	Hosts   []*Host
	History History
	Config  RunnerConfig
}

func NewRunner(hosts []*Host) *Runner {
	return &Runner{
		Hosts:   hosts,
		History: make([]HistoryItem, 0),
		Config: RunnerConfig{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *Runner) Run(command string) HistoryItem {
	hi := r.NewHistoryItem(command)
	c := make(chan Result)
	for _, host := range hi.Hosts {
		go host.Run(command, c)
	}
	timeout := time.After(r.Config.Timeout)
	todo := len(hi.Hosts)
	for todo > 0 {
		select {
		case <-timeout:
			r.Cancel(&hi)
		case r := <-c:
			hi.Results[r.Host] = r
			todo -= 1
		}
	}
	r.History = append(r.History, hi)
	return hi
}

func (r *Runner) NewHistoryItem(command string) HistoryItem {
	return HistoryItem{
		Hosts:   r.Hosts,
		Command: command,
		Results: make(map[string]Result),
	}
}

func (r *Runner) Cancel(hi *HistoryItem) {
	// Cancel nonfinished SSH runs
	// Add artificial timeout results for unfinished hosts
}
