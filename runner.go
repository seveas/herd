package herd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Result struct {
	Host       string
	ExitStatus int
	Stdout     []byte
	Stderr     []byte
	Err        error
	StartTime  time.Time
	EndTime    time.Time
}

func (r Result) MarshalJSON() ([]byte, error) {
	r_ := map[string]interface{}{
		"Host":       r.Host,
		"ExitStatus": r.ExitStatus,
		"Stdout":     string(r.Stdout),
		"Stderr":     string(r.Stderr),
		"Err":        r.Err,
		"ErrString":  "",
		"StartTime":  r.StartTime,
		"EndTime":    r.EndTime,
	}
	if r.Err != nil {
		r_["ErrString"] = r.Err.Error()
	}
	return json.Marshal(r_)
}

func (r Result) String() string {
	return fmt.Sprintf("[%s] (Err: %s)]\n%s\n---\n%s\n", r.Host, r.Err, string(r.Stdout), string(r.Stderr))
}

type HistoryItem struct {
	Hosts     []*Host
	Command   string
	Results   map[string]Result
	StartTime time.Time
	EndTime   time.Time
}

type History []HistoryItem

func (h History) Save(path string) {
	data, err := json.Marshal(h)
	if err != nil {
		UI.Warnf("Unable to export history: %s", err)
		return
	}
	err = ioutil.WriteFile(path, data, 0600)
	if err != nil {
		UI.Warnf("Unable to save history to %s: %s", path, err)
	}
}

type Runner struct {
	Hosts   []*Host
	History History
	Config  RunnerConfig
}

func NewRunner(hosts []*Host, config RunnerConfig) *Runner {
	return &Runner{
		Hosts:   hosts,
		History: make([]HistoryItem, 0),
		Config:  config,
	}
}

func (r *Runner) Run(command string) HistoryItem {
	hi := r.NewHistoryItem(command)
	c := make(chan Result)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, host := range hi.Hosts {
		hctx, hcancel := context.WithTimeout(ctx, r.Config.HostTimeout)
		host.SshConfig.Timeout = r.Config.ConnectTimeout
		go func(ctx context.Context, h *Host) { defer hcancel(); h.Run(ctx, command, c) }(hctx, host)
	}
	ticker := time.NewTicker(time.Second / 2)
	timeout := time.After(r.Config.Timeout)
	todo := len(hi.Hosts)
	total, doneOk, doneFail, doneError := todo, 0, 0, 0
	for todo > 0 {
		select {
		case <-ticker.C:
			UI.Progress(total, doneOk, doneFail, doneError, todo)
		case <-timeout:
			UI.Errorf("\nRun canceled with %d unfinished tasks!", todo)
			cancel()
		case r := <-c:
			if r.ExitStatus == -1 {
				doneError++
			} else if r.ExitStatus == 0 {
				doneOk++
			} else {
				doneFail++
			}
			hi.Results[r.Host] = r
			todo--
		}
		UI.Progress(total, doneOk, doneFail, doneError, todo)
	}
	hi.EndTime = time.Now()
	ticker.Stop()
	r.History = append(r.History, hi)
	return hi
}

func (r *Runner) NewHistoryItem(command string) HistoryItem {
	return HistoryItem{
		Hosts:     r.Hosts,
		Command:   command,
		Results:   make(map[string]Result),
		StartTime: time.Now(),
	}
}
