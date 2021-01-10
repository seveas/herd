package herd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/seveas/scattergather"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
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
	registry        *Registry
	hosts           Hosts
	parallel        int
	splay           time.Duration
	timeout         time.Duration
	hostTimeout     time.Duration
	connectTimeout  time.Duration
	sshAgentTimeout time.Duration
	agent           agent.Agent
	signers         []ssh.Signer
	signersByPath   map[string]ssh.Signer
}

type ProgressState int

const (
	Queued = iota
	Waiting
	Running
	Finished
)

func NewRunner(registry *Registry) *Runner {
	return &Runner{
		hosts:           make(Hosts, 0),
		registry:        registry,
		timeout:         60 * time.Second,
		hostTimeout:     10 * time.Second,
		connectTimeout:  3 * time.Second,
		sshAgentTimeout: 50 * time.Millisecond,
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

func (r *Runner) AddHosts(glob string, attrs MatchAttributes) {
	hosts := append(r.hosts, r.registry.GetHosts(glob, attrs)...)
	if !strings.HasPrefix(glob, "file:") {
		hosts.Sort(r.registry.sort)
	}
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
	if r.agent == nil && !strings.HasPrefix(command, "herd:keyscan:") {
		sock, err := agentConnection()
		if err != nil {
			logrus.Errorf("Unable to connect to ssh agent: %s", err)
			return hi
		}
		if _, ok := sock.(*net.UnixConn); ok {
			// Determine whether we can use the faster pipelined ssh agent protocol
			sock2, _ := agentConnection()
			sock3, _ := agentConnection()
			fastAgent := NewSshAgentClient(sock2, sock3)
			if fastAgent.functional(r.sshAgentTimeout) {
				sock.(*net.UnixConn).Close()
				r.agent = fastAgent
			} else {
				// Pity.
				logrus.Warnf("Using slow ssh agent, see https://herd.seveas.net/documentation/ssh_agent.html to fix this")
				sock2.(*net.UnixConn).Close()
				sock3.(*net.UnixConn).Close()
				r.agent = agent.NewClient(sock)
			}
		} else {
			r.agent = agent.NewClient(sock)
		}

		r.signers, err = r.agent.Signers()
		r.signersByPath = make(map[string]ssh.Signer)
		if err != nil {
			logrus.Errorf("Unable to retrieve keys from SSH agent: %s", err)
			return hi
		}
		if len(r.signers) == 0 {
			logrus.Errorf("No keys found in ssh agent")
			return hi
		}

		for _, signer := range r.signers {
			comment := signer.PublicKey().(*agent.Key).Comment
			r.signersByPath[comment] = signer
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var sg *scattergather.ScatterGather
	if r.parallel > 0 {
		sg = scattergather.New(int64(r.parallel))
	} else {
		sg = scattergather.New(int64(len(r.hosts)))
	}
	for _, host := range hi.Hosts {
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
			if len(host.sshConfig.Auth) == 0 {
				host.sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeysCallback(func() ([]ssh.Signer, error) { return r.signersForHost(host), nil })}
			}
			result := host.Run(ctx, command, oc)
			pc <- ProgressMessage{Host: host, State: Finished, Result: result}
			return result, nil
		}, ctx, host)
	}
	go func() {
		timeout := time.After(r.timeout)
		signals := make(chan os.Signal)
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
	}
	for _, host := range r.hosts {
		if _, ok := hi.Results[host.Name]; !ok {
			result := &Result{Host: host, ExitStatus: -1, Err: errors.New("context canceled")}
			host.lastResult = result
			pc <- ProgressMessage{Host: host, State: Finished, Result: result}
			hi.Results[host.Name] = result
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

func (r *Runner) signersForHost(h *Host) []ssh.Signer {
	if path := h.extConfig["identityfile"]; path != "" {
		if k, ok := r.signersByPath[path]; ok {
			return []ssh.Signer{k}
		} else {
			return []ssh.Signer{}
		}
	}
	return r.signers
}
