package herd

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"os/signal"
	"sort"
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
	hosts          *HostSet
	parallel       int
	splay          time.Duration
	timeout        time.Duration
	hostTimeout    time.Duration
	executor       Executor
	current        *scattergather.ScatterGather[*Result]
	cancel         context.CancelFunc
	signalHandlers map[os.Signal]func()
}

type ProgressState int

const (
	Queued ProgressState = iota
	Waiting
	Running
	Finished
)

func NewRunner(hosts *HostSet, executor Executor) *Runner {
	return &Runner{
		hosts:          hosts,
		executor:       executor,
		signalHandlers: make(map[os.Signal]func()),
	}
}

func (r *Runner) SetParallel(p int) {
	r.parallel = p
	if r.current != nil {
		r.current.SetParallel(int64(p))
	}
}

func (r *Runner) SetSplay(t time.Duration) {
	r.splay = t
}

func (r *Runner) SetTimeout(t time.Duration) {
	r.timeout = t
}

func (r *Runner) GetTimeout() time.Duration {
	if r.timeout != 0 {
		return r.timeout
	}
	if r.parallel == 0 || r.hostTimeout == 0 {
		return r.hostTimeout
	}
	batches := len(r.hosts.hosts)/r.parallel + 1
	return r.hostTimeout * time.Duration(batches)
}

func (r *Runner) SetHostTimeout(t time.Duration) {
	r.hostTimeout = t
}

func (r *Runner) GetHostTimeout() time.Duration {
	if r.hostTimeout != 0 {
		return r.hostTimeout
	}
	if r.parallel == 0 || r.timeout == 0 {
		return r.timeout
	}
	batches := len(r.hosts.hosts)/r.parallel + 1
	return r.timeout / time.Duration(batches)
}

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

func (r *Runner) Run(command string, pc chan ProgressMessage, oc chan OutputLine) (*HistoryItem, error) {
	if r.executor == nil {
		return nil, errors.New("No executor defined")
	}
	if len(r.hosts.hosts) == 0 {
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
	hi := newHistoryItem(command, len(r.hosts.hosts))
	hi.maxHostNameLength = r.hosts.maxNameLength
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	defer cancel()
	count := r.parallel
	if count <= 0 {
		count = len(r.hosts.hosts)
	}
	r.current = scattergather.New[*Result](int64(count))
	for index, host := range r.hosts.hosts {
		r.current.Run(ctx, func() (*Result, error) {
			if r.splay > 0 {
				pc <- ProgressMessage{Host: host, State: Waiting}
				r.splayDelay(ctx)
			}
			pc <- ProgressMessage{Host: host, State: Running}
			ctx, cancel := context.WithTimeout(ctx, r.GetHostTimeout())
			defer cancel()
			result := r.executor.Run(ctx, host, command, oc)
			result.index = index
			host.lastResult = result
			pc <- ProgressMessage{Host: host, State: Finished, Result: result}
			return result, nil
		})
	}
	go func() {
		timeout := time.After(r.GetTimeout())
		signals := make(chan os.Signal, 5)
		for s := range r.signalHandlers {
			signal.Notify(signals, s)
		}
		defer func() {
			for s := range r.signalHandlers {
				signal.Reset(s)
			}
		}()
		for {
			select {
			case <-timeout:
				logrus.Errorf("Run timed out with unfinished tasks!")
				cancel()
			case s := <-signals:
				r.signalHandlers[s]()
			case <-ctx.Done():
				return
			}
		}
	}()

	results, _ := r.current.Wait()
	r.current = nil
	r.cancel = nil
	cancel()
	for _, result := range results {
		hi.Results[result.index] = result
		switch result.ExitStatus {
		case -1:
			hi.Summary.Err++
		case 0:
			hi.Summary.Ok++
		default:
			hi.Summary.Fail++
		}
	}
	for index, host := range r.hosts.hosts {
		if hi.Results[index] == nil {
			result := &Result{Host: host.Name, ExitStatus: -1, Err: errors.New("context canceled")}
			host.lastResult = result
			pc <- ProgressMessage{Host: host, State: Finished, Result: result}
			hi.Results[index] = result
			hi.Summary.Err++
		}
	}
	for _, key := range r.hosts.sort {
		if key == "stdout" || key == "stderr" || key == "exitstatus" {
			// We re-sort hosts and results according to the result of the last command
			r.hosts.Sort()
			for idx, host := range r.hosts.hosts {
				host.lastResult.index = idx
			}
			sort.Slice(hi.Results, func(i, j int) bool { return hi.Results[i].index < hi.Results[j].index })
			break
		}
	}
	hi.end()
	return hi, nil
}

func (r *Runner) OnSignal(s os.Signal, f func()) {
	r.signalHandlers[s] = f
}

func (r *Runner) Interrupt() {
	if r.cancel != nil {
		r.cancel()
	}
}

func (r *Runner) End() {
	for _, h := range r.hosts.hosts {
		if h.Connection != nil {
			logrus.Debugf("Disconnecting from %s", h.Name)
			h.Connection.Close()
			h.Connection = nil
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
