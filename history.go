package katyusha

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type History []*HistoryItem

type HistoryItem struct {
	Hosts       Hosts
	Command     string
	Results     map[string]Result
	StartTime   time.Time
	EndTime     time.Time
	ElapsedTime float64
}

type Result struct {
	Host        string
	ExitStatus  int
	Stdout      []byte
	Stderr      []byte
	Err         error
	StartTime   time.Time
	EndTime     time.Time
	ElapsedTime float64
}

func NewHistoryItem(command string, hosts Hosts) *HistoryItem {
	return &HistoryItem{
		Hosts:     hosts,
		Command:   command,
		Results:   make(map[string]Result),
		StartTime: time.Now(),
	}
}

func (h *HistoryItem) End() {
	h.EndTime = time.Now()
	h.ElapsedTime = h.EndTime.Sub(h.StartTime).Seconds()
}

func (r Result) MarshalJSON() ([]byte, error) {
	r_ := map[string]interface{}{
		"Host":        r.Host,
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
	data, err := json.Marshal(h)
	if err != nil {
		UI.Warnf("Unable to export history: %s", err)
		return err
	}
	err = ioutil.WriteFile(path, data, 0600)
	if err != nil {
		UI.Warnf("Unable to save history to %s: %s", path, err)
	}
	return err
}
