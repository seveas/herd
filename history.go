package herd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

type History []*HistoryItem

type HistoryItem struct {
	Command string
	Results []*Result
	Summary struct {
		Ok   int
		Fail int
		Err  int
	}
	StartTime         time.Time
	EndTime           time.Time
	ElapsedTime       float64
	maxHostNameLength int
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
	index       int
}

type resultx struct {
	Host        string
	ExitStatus  int
	Stdout      string
	Stderr      string
	Err         any
	ErrString   string
	StartTime   time.Time
	EndTime     time.Time
	ElapsedTime float64
}

func newHistoryItem(command string, nhosts int) *HistoryItem {
	return &HistoryItem{
		Command:   command,
		Results:   make([]*Result, nhosts),
		StartTime: time.Now(),
	}
}

func (h *HistoryItem) MarshalJSON() ([]byte, error) {
	r := map[string]interface{}{
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
	r_ := resultx{
		Host:        r.Host,
		ExitStatus:  r.ExitStatus,
		Stdout:      string(r.Stdout),
		Stderr:      string(r.Stderr),
		Err:         r.Err,
		ErrString:   "",
		StartTime:   r.StartTime,
		EndTime:     r.EndTime,
		ElapsedTime: r.ElapsedTime,
	}
	if r.Err != nil {
		r_.ErrString = r.Err.Error()
	}
	return json.Marshal(r_)
}

func (r *Result) UnmarshalJSON(data []byte) error {
	r_ := resultx{}
	if err := json.Unmarshal(data, &r_); err != nil {
		return err
	}
	r.Host = r_.Host
	r.ExitStatus = r_.ExitStatus
	r.Stdout = []byte(r_.Stdout)
	r.Stderr = []byte(r_.Stderr)
	r.StartTime = r_.StartTime
	r.EndTime = r_.EndTime
	r.ElapsedTime = r_.ElapsedTime
	return nil
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
	if err = os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		logrus.Warnf("Unable to create history path %s: %s", filepath.Dir(path), err)
		return err
	}
	if err = os.WriteFile(path, data, 0o600); err != nil {
		logrus.Warnf("Unable to save history to %s: %s", path, err)
	} else {
		logrus.Infof("History saved to %s", path)
	}
	return err
}
