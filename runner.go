package herd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
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
	Hosts     Hosts
	Providers Providers
	History   History
	Config    *RunnerConfig
}

func NewRunner(providers Providers, config *RunnerConfig) *Runner {
	return &Runner{
		Hosts:     make([]*Host, 0),
		Providers: providers,
		History:   make([]HistoryItem, 0),
		Config:    config,
	}
}

func (r *Runner) AddHosts(glob string, attrs HostAttributes) {
	hosts := append(r.Hosts, r.Providers.GetHosts(glob, attrs)...)
	r.Hosts = hosts.SortAndUniq()
}

func (r *Runner) RemoveHosts(glob string, attrs HostAttributes) {
	newHosts := make([]*Host, 0)
	for _, host := range r.Hosts {
		if !host.Match(glob, attrs) {
			newHosts = append(newHosts, host)
		}
	}
	r.Hosts = newHosts
}

func (r *Runner) ListHosts(oneline bool) {
	if oneline {
		hosts := strings.Builder{}
		for i, host := range r.Hosts {
			if i == 0 {
				fmt.Fprint(&hosts, host.Name)
			} else {
				fmt.Fprintf(&hosts, ",%s", host.Name)
			}
		}
		UI.Println(hosts.String())
	} else {
		for _, host := range r.Hosts {
			UI.Println(host.Name)
		}
	}
}

func (r *Runner) Run(command string) HistoryItem {
	hi := r.NewHistoryItem(command)
	c := make(chan Result)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queued := -1
	todo := len(hi.Hosts)
	total, doneOk, doneFail, doneError := todo, 0, 0, 0
	if r.Config.Parallel > 0 {
		queued = len(hi.Hosts)
		hqueue := make(chan *Host)
		go func() {
			for _, host := range hi.Hosts {
				hqueue <- host
				queued--
				UI.Progress(total, todo, queued, doneOk, doneFail, doneError)
			}
			close(hqueue)
		}()
		for i := 0; i < r.Config.Parallel; i++ {
			UI.Debugf("Starting worker %d/%d", i+1, r.Config.Parallel)
			go func() {
				for {
					host, ok := <-hqueue
					if !ok {
						break
					}
					host.SshConfig.Timeout = r.Config.ConnectTimeout
					hctx, hcancel := context.WithTimeout(ctx, r.Config.HostTimeout)
					host.Run(hctx, command, c)
					hcancel()
				}
			}()
		}
	} else {
		for _, host := range hi.Hosts {
			hctx, hcancel := context.WithTimeout(ctx, r.Config.HostTimeout)
			host.SshConfig.Timeout = r.Config.ConnectTimeout
			go func(ctx context.Context, h *Host) { defer hcancel(); h.Run(ctx, command, c) }(hctx, host)
		}
	}
	ticker := time.NewTicker(time.Second / 2)
	timeout := time.After(r.Config.Timeout)
	for todo > 0 {
		select {
		case <-ticker.C:
			UI.Progress(total, todo, queued, doneOk, doneFail, doneError)
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
		UI.Progress(total, todo, queued, doneOk, doneFail, doneError)
	}
	hi.EndTime = time.Now()
	ticker.Stop()
	r.History = append(r.History, hi)
	return hi
}

func (r *Runner) End() {
	for _, host := range r.Hosts {
		host.Disconnect()
	}
}

func (r *Runner) NewHistoryItem(command string) HistoryItem {
	return HistoryItem{
		Hosts:     r.Hosts,
		Command:   command,
		Results:   make(map[string]Result),
		StartTime: time.Now(),
	}
}
