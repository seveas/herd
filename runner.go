package herd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"os/signal"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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

func (r *Runner) ListHosts(oneline, allAttributes bool, attributes []string, csvOutput bool) {
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
	} else if allAttributes || len(attributes) > 0 {
		var writer datawriter
		if csvOutput {
			writer = csv.NewWriter(UI)
		} else {
			writer = NewColumnizer(UI, "   ")
		}
		if allAttributes {
			attrs := make(map[string]bool)
			for _, host := range r.Hosts {
				for key, _ := range host.Attributes {
					attrs[key] = true
				}
			}
			for attr, _ := range attrs {
				attributes = append(attributes, attr)
			}
			sort.Strings(attributes)
			attrline := make([]string, len(attributes)+1)
			attrline[0] = "name"
			copy(attrline[1:], attributes)
			writer.Write(attrline)
		}
		for _, host := range r.Hosts {
			line := make([]string, len(attributes)+1)
			line[0] = host.Name
			for i, attr := range attributes {
				val, ok := host.Attributes[attr]
				if ok {
					line[i+1] = fmt.Sprintf("%v", val)
				} else {
					line[i+1] = ""
				}
			}
			writer.Write(line)
		}
		writer.Flush()
	} else {
		for _, host := range r.Hosts {
			UI.Println(host.Name)
		}
	}
}

func (r *Runner) Run(command string) {
	hi := NewHistoryItem(command, r.Hosts)
	c := make(chan Result)
	defer close(c)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if viper.GetString("Output") == "line" {
		ctx = context.WithValue(ctx, "hostnamelen", maxHostNameLen(hi.Hosts))
		UI.PrintCommand(hi.Command)
	}
	queued := -1
	todo := len(hi.Hosts)
	total, doneOk, doneFail, doneError := todo, 0, 0, 0
	if viper.GetInt("Parallel") > 0 {
		queued = len(hi.Hosts)
		hqueue := make(chan *Host)
		go func() {
			for _, host := range hi.Hosts {
				hqueue <- host
				queued--
				UI.Progress(hi.StartTime, total, todo, queued, doneOk, doneFail, doneError)
			}
			close(hqueue)
		}()
		for i := 0; i < viper.GetInt("Parallel"); i++ {
			logrus.Debugf("Starting worker %d/%d", i+1, viper.GetInt("Parallel"))
			go func() {
				for host := range hqueue {
					host.SshConfig.Timeout = viper.GetDuration("ConnectTimeout")
					hctx, hcancel := context.WithTimeout(ctx, viper.GetDuration("HostTimeout"))
					c <- host.Run(hctx, command)
					hcancel()
				}
			}()
		}
	} else {
		for _, host := range hi.Hosts {
			hctx, hcancel := context.WithTimeout(ctx, viper.GetDuration("HostTimeout"))
			host.SshConfig.Timeout = viper.GetDuration("ConnectTimeout")
			go func(ctx context.Context, h *Host) { c <- h.Run(ctx, command); hcancel() }(hctx, host)
		}
	}
	ticker := time.NewTicker(time.Second / 2)
	defer ticker.Stop()
	timeout := time.After(viper.GetDuration("Timeout"))
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	defer signal.Reset(os.Interrupt)
	for todo > 0 {
		select {
		case <-ticker.C:
			UI.Progress(hi.StartTime, total, todo, queued, doneOk, doneFail, doneError)
		case <-timeout:
			logrus.Errorf("Run canceled with %d unfinished tasks!", todo)
			cancel()
		case <-signals:
			logrus.Errorf("Interrupted, canceling with %d unfinished tasks", todo)
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
			if viper.GetString("Output") == "host" {
				UI.PrintResult(r)
			}
			todo--
		}
		UI.Progress(hi.StartTime, total, todo, queued, doneOk, doneFail, doneError)
	}
	hi.End()
	r.History = append(r.History, hi)
	if viper.GetString("Output") == "all" {
		UI.PrintHistoryItem(hi)
	}
	if viper.GetString("Output") == "pager" {
		UI.PrintHistoryItemWithPager(hi)
	}
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

func maxHostNameLen(hosts Hosts) int {
	max := 0
	for _, host := range hosts {
		if len(host.Name) > max {
			max = len(host.Name)
		}
	}
	return max
}
