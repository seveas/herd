package katyusha

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
)

type OutputLine struct {
	Host   *Host
	Stderr bool
	Data   []byte
}

type ProgressMessage struct {
	Host   *Host
	Result *Result
}

type Runner struct {
	registry       *Registry
	hosts          Hosts
	parallel       int
	splay          time.Duration
	timeout        time.Duration
	hostTimeout    time.Duration
	connectTimeout time.Duration
}

func NewRunner(registry *Registry) *Runner {
	return &Runner{
		hosts:          make(Hosts, 0),
		registry:       registry,
		timeout:        60 * time.Second,
		hostTimeout:    10 * time.Second,
		connectTimeout: 3 * time.Second,
	}
}

func (r *Runner) GetHosts() Hosts {
	return r.hosts[:]
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

func (r *Runner) AddHosts(glob string, attrs MatchAttributes) {
	hosts := append(r.hosts, r.registry.GetHosts(glob, attrs)...)
	hosts.Sort(r.registry.sort)
	r.hosts = hosts.Uniq()
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
	if pc == nil {
		pc = make(chan ProgressMessage)
		go func() {
			for range pc {
			}
		}()
	}
	defer close(pc)
	if oc != nil {
		defer close(oc)
	}
	hi := newHistoryItem(command, r.hosts)
	c := make(chan *Result)
	defer close(c)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if r.parallel > 0 {
		hqueue := make(chan *Host)
		go func() {
			for _, host := range hi.Hosts {
				hqueue <- host
			}
			close(hqueue)
		}()
		for i := 0; i < r.parallel; i++ {
			logrus.Debugf("Starting worker %d/%d", i+1, r.parallel)
			go func() {
				for host := range hqueue {
					pc <- ProgressMessage{Host: host}
					r.splayDelay(ctx)
					host.sshConfig.Timeout = r.connectTimeout
					hctx, hcancel := context.WithTimeout(ctx, r.hostTimeout)
					c <- host.Run(hctx, command, oc)
					hcancel()
				}
			}()
		}
	} else {
		for _, host := range hi.Hosts {
			go func(ctx context.Context, h *Host) {
				pc <- ProgressMessage{Host: h}
				r.splayDelay(ctx)
				ctx, hcancel := context.WithTimeout(ctx, r.hostTimeout)
				h.sshConfig.Timeout = r.connectTimeout
				c <- h.Run(ctx, command, oc)
				hcancel()
			}(ctx, host)
		}
	}
	timeout := time.After(r.timeout)
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	defer signal.Reset(os.Interrupt)
	todo := len(r.hosts)
	for todo > 0 {
		select {
		case <-timeout:
			logrus.Errorf("Run canceled with %d unfinished tasks!", todo)
			cancel()
		case <-signals:
			logrus.Errorf("Interrupted, canceling with %d unfinished tasks", todo)
			cancel()
		case r := <-c:
			pc <- ProgressMessage{Host: r.Host, Result: r}
			hi.Results[r.Host.Name] = r
			todo--
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
