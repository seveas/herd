package herd

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/seveas/scattergather"
	"github.com/sirupsen/logrus"
)

type Executor interface {
	Run(ctx context.Context, host *Host, cmd string, oc chan OutputLine) *Result
	SetConnectTimeout(time.Duration)
}

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
	hosts       Hosts
	sort        []string
	parallel    int
	splay       time.Duration
	timeout     time.Duration
	hostTimeout time.Duration
	executor    Executor
}

type ProgressState int

const (
	Queued = iota
	Waiting
	Running
	Finished
)

func NewRunner(executor Executor) *Runner {
	return &Runner{
		hosts:       make(Hosts, 0),
		sort:        []string{"name"},
		executor:    executor,
		timeout:     60 * time.Second,
		hostTimeout: 10 * time.Second,
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

// FIXME
func (r *Runner) SetConnectTimeout(t time.Duration) {
	if r.executor != nil {
		r.executor.SetConnectTimeout(t)
	}
}

func (r *Runner) Settings() (string, map[string]interface{}) {
	return "Runner", map[string]interface{}{
		"Parallel":    r.parallel,
		"Splay":       r.splay,
		"Timeout":     r.timeout,
		"HostTimeout": r.hostTimeout,
	}
}

func (r *Runner) AddHosts(hosts Hosts) {
	h := r.hosts[:]
	h = append(h, hosts...)
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

func (r *Runner) Run(command string, pc chan ProgressMessage, oc chan OutputLine) (*HistoryItem, error) {
	if r.executor == nil {
		return nil, errors.New("No executor defined")
	}
	if len(r.hosts) == 0 {
		return nil, errors.New("No hosts selected")
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
	count := r.parallel
	if count <= 0 {
		count = len(r.hosts)
	}
	sg := scattergather.New[*Result](int64(count))
	for _, host := range hi.Hosts {
		host := host
		sg.Run(ctx, func() (*Result, error) {
			if r.splay > 0 {
				pc <- ProgressMessage{Host: host, State: Waiting}
				r.splayDelay(ctx)
			}
			pc <- ProgressMessage{Host: host, State: Running}
			ctx, cancel := context.WithTimeout(ctx, r.hostTimeout)
			defer cancel()
			result := r.executor.Run(ctx, host, command, oc)
			host.lastResult = result
			pc <- ProgressMessage{Host: host, State: Finished, Result: result}
			return result, nil
		})
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
	for _, result := range results {
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
	return hi, nil
}

func (r *Runner) End() {
	for _, h := range r.hosts {
		if h.Connection != nil {
			logrus.Debugf("Disconnecting from %s", h.Name)
			h.Connection.Close()
		}
	}
}

func (r *Runner) splayDelay(ctx context.Context) {
	if r.splay <= 0 {
		return
	}
	d := time.Duration(rand.Int63n(int64(r.splay))) //#nosec G404 -- This does not need cryptographically secure numbers
	tctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()
	select {
	case <-ctx.Done():
		return
	case <-tctx.Done():
		return
	}
}

type TimeoutError struct {
	Message string
}

func (e TimeoutError) Error() string {
	return e.Message
}
