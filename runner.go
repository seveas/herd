package katyusha

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

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
	registry       *Registry
	hosts          Hosts
	parallel       int
	splay          time.Duration
	timeout        time.Duration
	hostTimeout    time.Duration
	connectTimeout time.Duration
	agent          agent.Agent
	signers        []ssh.Signer
	signersByPath  map[string]ssh.Signer
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
	if r.agent == nil {
		sock, err := agentConnection()
		if err != nil {
			logrus.Errorf("Unable to connect to ssh agent: %s", err)
			return hi
		}
		if _, ok := sock.(*net.UnixConn); ok {
			if _, ok := os.LookupEnv("KATYUSHA_FAST_SSH_AGENT"); ok {
				sock2, _ := agentConnection()
				r.agent = NewSshAgentClient(sock, sock2)
			} else {
				if _, ok := os.LookupEnv("KATYUSHA_SLOW_SSH_AGENT"); !ok {
					logrus.Warnf("Using slow ssh agent, see https://katyusha.seveas.net/documentation/ssh_agent.html to fix this")
				}
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
					if r.splay > 0 {
						pc <- ProgressMessage{Host: host, State: Waiting}
						r.splayDelay(ctx)
					}
					pc <- ProgressMessage{Host: host, State: Running}
					host.sshConfig.Timeout = r.connectTimeout
					if len(host.sshConfig.Auth) == 0 {
						host.sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeysCallback(func() ([]ssh.Signer, error) { return r.signersForHost(host), nil })}
					}
					hctx, hcancel := context.WithTimeout(ctx, r.hostTimeout)
					c <- host.Run(hctx, command, oc)
					hcancel()
				}
			}()
		}
	} else {
		for _, host := range hi.Hosts {
			go func(ctx context.Context, h *Host) {
				if r.splay > 0 {
					pc <- ProgressMessage{Host: h, State: Waiting}
					r.splayDelay(ctx)
				}
				pc <- ProgressMessage{Host: h, State: Running}
				ctx, hcancel := context.WithTimeout(ctx, r.hostTimeout)
				h.sshConfig.Timeout = r.connectTimeout
				if len(h.sshConfig.Auth) == 0 {
					h.sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeysCallback(func() ([]ssh.Signer, error) { return r.signersForHost(h), nil })}
				}
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
			pc <- ProgressMessage{Host: r.Host, State: Finished, Result: r}
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
