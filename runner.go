package herd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"

	"github.com/spf13/viper"
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
}

func NewRunner(providers Providers) *Runner {
	return &Runner{
		Hosts:     make([]*Host, 0),
		Providers: providers,
		History:   make([]HistoryItem, 0),
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
	defer close(c)
	ctx, cancel := context.WithCancel(context.Background())
	if viper.GetString("Output") == "line" {
		ctx = context.WithValue(ctx, "hostnamelen", maxHostNameLen(hi.Hosts))
	}
	defer cancel()
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
				UI.Progress(total, todo, queued, doneOk, doneFail, doneError)
			}
			close(hqueue)
		}()
		for i := 0; i < viper.GetInt("Parallel"); i++ {
			UI.Debugf("Starting worker %d/%d", i+1, viper.GetInt("Parallel"))
			go func() {
				for host := range hqueue {
					host.SshConfig.Timeout = viper.GetDuration("ConnectTimeout")
					hctx, hcancel := context.WithTimeout(ctx, viper.GetDuration("HostTimeout"))
					host.Run(hctx, command, c)
					hcancel()
				}
			}()
		}
	} else {
		for _, host := range hi.Hosts {
			hctx, hcancel := context.WithTimeout(ctx, viper.GetDuration("HostTimeout"))
			host.SshConfig.Timeout = viper.GetDuration("ConnectTimeout")
			go func(ctx context.Context, h *Host) { defer hcancel(); h.Run(ctx, command, c) }(hctx, host)
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
			UI.Progress(total, todo, queued, doneOk, doneFail, doneError)
		case <-timeout:
			UI.Errorf("Run canceled with %d unfinished tasks!", todo)
			cancel()
		case <-signals:
			UI.Errorf("Interrupted, canceling with %d unfinished tasks", todo)
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
		UI.Progress(total, todo, queued, doneOk, doneFail, doneError)
	}
	hi.EndTime = time.Now()
	r.History = append(r.History, hi)
	if viper.GetString("Output") == "all" {
		UI.PrintHistoryItem(hi)
	}
	return hi
}

func (r *Runner) End() {
	for _, host := range r.Hosts {
		host.Disconnect()
	}

	// Save history, if there is any
	if len(r.History) > 0 {
		if err := os.MkdirAll(viper.GetString("HistoryDir"), 0700); err != nil {
			UI.Warnf("Unable to create history path %s: %s", viper.GetString("HistoryDir"), err)
		} else {
			fn := path.Join(viper.GetString("HistoryDir"), r.History[0].StartTime.Format("2006-01-02T15:04:05.json"))
			r.History.Save(fn)
		}
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

func maxHostNameLen(hosts []*Host) int {
	max := 0
	for _, host := range hosts {
		if len(host.Name) > max {
			max = len(host.Name)
		}
	}
	return max
}
