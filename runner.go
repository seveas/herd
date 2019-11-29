package herd

import (
	"context"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	Registry *Registry
	Hosts    Hosts
	History  History
}

func NewRunner(registry *Registry) *Runner {
	return &Runner{
		Hosts:    make(Hosts, 0),
		Registry: registry,
		History:  make(History, 0),
	}
}

func (r *Runner) AddHosts(glob string, attrs MatchAttributes) {
	hosts := append(r.Hosts, r.Registry.GetHosts(glob, attrs)...)
	hosts.Sort()
	r.Hosts = hosts.Uniq()
}

func (r *Runner) RemoveHosts(glob string, attrs MatchAttributes) {
	newHosts := make(Hosts, 0)
	for _, host := range r.Hosts {
		if !host.Match(glob, attrs) {
			newHosts = append(newHosts, host)
		}
	}
	r.Hosts = newHosts
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
	hi := NewHistoryItem(command, r.Hosts)
	c := make(chan *Result)
	defer close(c)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if viper.GetInt("Parallel") > 0 {
		hqueue := make(chan *Host)
		go func() {
			for _, host := range hi.Hosts {
				hqueue <- host
			}
			close(hqueue)
		}()
		for i := 0; i < viper.GetInt("Parallel"); i++ {
			logrus.Debugf("Starting worker %d/%d", i+1, viper.GetInt("Parallel"))
			go func() {
				for host := range hqueue {
					host.SshConfig.Timeout = viper.GetDuration("ConnectTimeout")
					hctx, hcancel := context.WithTimeout(ctx, viper.GetDuration("HostTimeout"))
					pc <- ProgressMessage{Host: host}
					c <- host.Run(hctx, command, oc)
					hcancel()
				}
			}()
		}
	} else {
		for _, host := range hi.Hosts {
			hctx, hcancel := context.WithTimeout(ctx, viper.GetDuration("HostTimeout"))
			host.SshConfig.Timeout = viper.GetDuration("ConnectTimeout")
			go func(ctx context.Context, h *Host) {
				pc <- ProgressMessage{Host: h}
				c <- h.Run(hctx, command, oc)
				hcancel()
			}(hctx, host)
		}
	}
	timeout := time.After(viper.GetDuration("Timeout"))
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	defer signal.Reset(os.Interrupt)
	todo := len(r.Hosts)
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
	hi.End()
	r.History = append(r.History, hi)
	return hi
}

func (r *Runner) End() error {
	for _, host := range r.Hosts {
		host.Disconnect()
	}

	// Save history, if there is any
	var err error
	if len(r.History) > 0 {
		if err = os.MkdirAll(viper.GetString("HistoryDir"), 0700); err != nil {
			logrus.Warnf("Unable to create history path %s: %s", viper.GetString("HistoryDir"), err)
		} else {
			fn := path.Join(viper.GetString("HistoryDir"), r.History[0].StartTime.Format("2006-01-02T15:04:05.json"))
			err = r.History.Save(fn)
			if err == nil {
				logrus.Infof("History saved to %s", fn)
			}
		}
	}
	return err
}
