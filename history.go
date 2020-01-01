package katyusha

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

type History []*HistoryItem

type HistoryItem struct {
	Hosts       Hosts
	Command     string
	Results     map[string]*Result
	StartTime   time.Time
	EndTime     time.Time
	ElapsedTime float64
}

type Result struct {
	Host        *Host
	ExitStatus  int
	Stdout      []byte
	Stderr      []byte
	Err         error
	StartTime   time.Time
	EndTime     time.Time
	ElapsedTime float64
}

func newHistoryItem(command string, hosts Hosts) *HistoryItem {
	return &HistoryItem{
		Hosts:     hosts,
		Command:   command,
		Results:   make(map[string]*Result),
		StartTime: time.Now(),
	}
}

func (h *HistoryItem) MarshalJSON() ([]byte, error) {
	hosts := make([]string, len(h.Hosts))
	for i, h_ := range h.Hosts {
		hosts[i] = h_.Name
	}
	r := map[string]interface{}{
		"Hosts":       hosts,
		"Command":     h.Command,
		"Results":     h.Results,
		"StartTime":   h.StartTime,
		"EndTime":     h.EndTime,
		"ElapsedTime": h.ElapsedTime,
	}
	return json.Marshal(r)
}

func (h *HistoryItem) end() {
	h.EndTime = time.Now()
	h.ElapsedTime = h.EndTime.Sub(h.StartTime).Seconds()
}

func (r Result) MarshalJSON() ([]byte, error) {
	r_ := map[string]interface{}{
		"Host":        r.Host.Name,
		"ExitStatus":  r.ExitStatus,
		"Stdout":      string(r.Stdout),
		"Stderr":      string(r.Stderr),
		"Err":         r.Err,
		"ErrString":   "",
		"StartTime":   r.StartTime,
		"EndTime":     r.EndTime,
		"ElapsedTime": r.ElapsedTime,
	}
	if r.Err != nil {
		r_["ErrString"] = r.Err.Error()
	}
	return json.Marshal(r_)
}

func (r Result) String() string {
	return fmt.Sprintf("[%s] (Err: %s)]\n%s\n---\n%s\n", r.Host, r.Err, string(r.Stdout), string(r.Stderr))
}

func (h History) Save(path string) error {
	if len(h) == 0 {
		return nil
	}
	data, err := json.Marshal(h)
	if err != nil {
		logrus.Warnf("Unable to export history: %s", err)
		return err
	}
	if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		logrus.Warnf("Unable to create history path %s: %s", filepath.Dir(path), err)
		return err
	}
	if err = ioutil.WriteFile(path, data, 0600); err != nil {
		logrus.Warnf("Unable to save history to %s: %s", path, err)
	} else {
		logrus.Infof("History saved to %s", path)
	}
	return err
}
