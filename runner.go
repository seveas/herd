package herd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/seveas/herd/sshagent"

	"github.com/seveas/scattergather"
	"github.com/sirupsen/logrus"
)

type OutputLine struct {
	Host   *Host
	Stderr bool
	Data   []byte
}

type ProgressMessage struct {
	Host   *Host
	State  ProgressState
	Result *Result
}

type Runner struct {
	hosts           Hosts
	sort            []string
	parallel        int
	splay           time.Duration
	timeout         time.Duration
	hostTimeout     time.Duration
	connectTimeout  time.Duration
	sshAgentTimeout time.Duration
	sshAgent        *sshagent.Agent
}

type ProgressState int

const (
	Queued = iota
	Waiting
	Running
	Finished
)

func NewRunner(agent *sshagent.Agent) *Runner {
	return &Runner{
		hosts:           make(Hosts, 0),
		sort:            []string{"name"},
		timeout:         60 * time.Second,
		hostTimeout:     10 * time.Second,
		connectTimeout:  3 * time.Second,
		sshAgent:        agent,
		sshAgentTimeout: 50 * time.Millisecond,
	}
}

func (r *Runner) GetHosts() Hosts {
	return r.hosts[:]
}

func (r *Runner) SetSortFields(s []string) {
	r.sort = s
}

func (r *Runner) SetParallel(p int) {
	r.parallel = p
}

func (r *Runner) SetSplay(t time.Duration) {
	r.splay = t
}

func (r *Runner) SetTimeout(t time.Duration) {
	r.timeout = t
}

func (r *Runner) SetHostTimeout(t time.Duration) {
	r.hostTimeout = t
}

func (r *Runner) SetConnectTimeout(t time.Duration) {
	r.connectTimeout = t
}

func (r *Runner) SetSshAgentTimeout(t time.Duration) {
	r.sshAgentTimeout = t
}

func (r *Runner) PrintSettings(ui io.Writer) {
	fmt.Fprintf(ui, "Parallel:       %d\n", r.parallel)
	fmt.Fprintf(ui, "Splay:          %s\n", r.splay)
	fmt.Fprintf(ui, "Timeout:        %s\n", r.timeout)
	fmt.Fprintf(ui, "HostTimeout:    %s\n", r.hostTimeout)
	fmt.Fprintf(ui, "ConnectTimeout: %s\n", r.connectTimeout)
}

func (r *Runner) AddHosts(hosts Hosts) {
	h := append(r.hosts, hosts...)
	h.Sort(r.sort)
	r.hosts = h.Uniq()
}

func (r *Runner) RemoveHosts(glob string, attrs MatchAttributes) {
	newHosts := make(Hosts, 0)
	for _, host := range r.hosts {
		if !host.Match(glob, attrs) {
			newHosts = append(newHosts, host)
		}
	}
	r.hosts = newHosts
}

func (r *Runner) Run(command string, pc chan ProgressMessage, oc chan OutputLine) *HistoryItem {
	if len(r.hosts) == 0 {
		logrus.Errorf("No hosts selected")
		return nil
	}
	if pc == nil {
		pc = make(chan ProgressMessage)
		defer close(pc)
		go func() {
			for range pc {
			}
		}()
	}
	hi := newHistoryItem(command, r.hosts)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var sg *scattergather.ScatterGather
	if r.parallel > 0 {
		sg = scattergather.New(int64(r.parallel))
	} else {
		sg = scattergather.New(int64(len(r.hosts)))
	}
	for _, host := range hi.Hosts {
		host.sshAgent = r.sshAgent
		sg.Run(func(ctx context.Context, args ...interface{}) (interface{}, error) {
			host := args[0].(*Host)
			if r.splay > 0 {
				pc <- ProgressMessage{Host: host, State: Waiting}
				r.splayDelay(ctx)
			}
			pc <- ProgressMessage{Host: host, State: Running}
			ctx, cancel := context.WithTimeout(ctx, r.hostTimeout)
			defer cancel()
			host.sshConfig.Timeout = r.connectTimeout
			result := host.Run(ctx, command, oc)
			pc <- ProgressMessage{Host: host, State: Finished, Result: result}
			return result, nil
		}, ctx, host)
	}
	go func() {
		timeout := time.After(r.timeout)
		signals := make(chan os.Signal, 5)
		signal.Notify(signals, os.Interrupt)
		defer signal.Reset(os.Interrupt)
		select {
		case <-timeout:
			logrus.Errorf("Run canceled with unfinished tasks!")
			cancel()
		case <-signals:
			logrus.Errorf("Interrupted, canceling with unfinished tasks")
			cancel()
		case <-ctx.Done():
			return
		}
	}()
	results, _ := sg.Wait()
	cancel()
	for _, rawResult := range results {
		result := rawResult.(*Result)
		hi.Results[result.Host.Name] = result
		switch result.ExitStatus {
		case -1:
			hi.Summary.Err++
		case 0:
			hi.Summary.Ok++
		default:
			hi.Summary.Fail++
		}
	}
	for _, host := range r.hosts {
		if _, ok := hi.Results[host.Name]; !ok {
			result := &Result{Host: host, ExitStatus: -1, Err: errors.New("context canceled")}
			host.lastResult = result
			pc <- ProgressMessage{Host: host, State: Finished, Result: result}
			hi.Results[host.Name] = result
			hi.Summary.Err++
		}
	}
	for _, key := range r.sort {
		if key == "stdout" || key == "stderr" || key == "exitstatus" {
			hi.Hosts.Sort(r.sort)
			break
		}
	}
	hi.end()
	return hi
}

func (r *Runner) End() {
	for _, host := range r.hosts {
		host.disconnect()
	}
}

func (r *Runner) splayDelay(ctx context.Context) {
	if r.splay <= 0 {
		return
	}
	d := time.Duration(rand.Int63n(int64(r.splay)))
	tctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()
	select {
	case <-ctx.Done():
		return
	case <-tctx.Done():
		return
	}
}
