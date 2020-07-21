package herd

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

var results []*Result
var testformatter = prettyFormatter{
	colors: map[logrus.Level]string{
		logrus.WarnLevel:  "yellow",
		logrus.ErrorLevel: "red+b",
		logrus.DebugLevel: "black+h",
	},
}

func init() {
	i := 0
	count := func() int {
		i++
		return i
	}
	start := time.Date(2019, 12, 8, 20, 26, 0, 0, time.UTC)
	results = []*Result{
		{
			Host:        NewHost(fmt.Sprintf("test-host-%03d.example.com", count()), HostAttributes{}),
			Err:         fmt.Errorf("It's always DNS"),
			ExitStatus:  -1,
			Stdout:      []byte{},
			Stderr:      []byte{},
			StartTime:   start,
			EndTime:     start.Add(12 * time.Second),
			ElapsedTime: 12,
		},
		{
			Host:        NewHost(fmt.Sprintf("test-host-%03d.example.com", count()), HostAttributes{}),
			Stdout:      []byte("May the forks be with you\nAnd you\n"),
			Stderr:      []byte{},
			StartTime:   start,
			EndTime:     start.Add(12 * time.Second),
			ElapsedTime: 12,
		},
		{
			Host:        NewHost(fmt.Sprintf("test-host-%03d.example.com", count()), HostAttributes{}),
			Stdout:      []byte("Newline is added automatically"),
			Stderr:      []byte{},
			StartTime:   start,
			EndTime:     start.Add(12 * time.Second),
			ElapsedTime: 12,
		},
		{
			Host:        NewHost(fmt.Sprintf("test-host-%03d.example.com", count()), HostAttributes{}),
			Err:         fmt.Errorf("Process exited with status 1"),
			ExitStatus:  1,
			Stdout:      []byte{},
			Stderr:      []byte("Text on stderr\nMore text\n"),
			StartTime:   start,
			EndTime:     start.Add(12 * time.Second),
			ElapsedTime: 12,
		},
		{
			Host:        NewHost(fmt.Sprintf("test-host-%03d.example.com", count()), HostAttributes{}),
			Stdout:      []byte("Text on stdout without newline"),
			Stderr:      []byte("Text on stderr\nMore text\n"),
			StartTime:   start,
			EndTime:     start.Add(12 * time.Second),
			ElapsedTime: 12,
		},
	}
}

func TestPrettyFormatterFormatCommand(t *testing.T) {
	if x := testformatter.formatCommand("hello world"); x != "\033[0;36mhello world\033[0m\n" {
		t.Errorf("Expected colored string, got %s", strconv.Quote(x))
	}
}

func TestPrettyFormatterFormatStatus(t *testing.T) {
	expected := []string{
		"\033[0;31mtest-host-001.example.com  It's always DNS after 12s\033[0m\n",
		"\033[0;32mtest-host-002.example.com  completed successfully after 12s\033[0m\n",
		"\033[0;32mtest-host-003.example.com  completed successfully after 12s\033[0m\n",
		"\033[0;31mtest-host-004.example.com  Process exited with status 1 after 12s\033[0m\n",
		"\033[0;32mtest-host-005.example.com  completed successfully after 12s\033[0m\n",
	}
	for i, r := range results {
		if s := testformatter.formatStatus(r, 0); s != expected[i] {
			t.Errorf("Result %d, expected status %s, got %s", i, strconv.Quote(expected[i]), strconv.Quote(s))
		}
	}
}

func TestPrettyFormatterFormatOutput(t *testing.T) {
	expected := []string{
		"\033[0;31mtest-host-001.example.com  It's always DNS after 12s\033[0m\n",
		"test-host-002.example.com  May the forks be with you\n  And you\n",
		"test-host-003.example.com  Newline is added automatically\n",
		"\033[0;31mtest-host-004.example.com  \033[0mText on stderr\n  More text\n\033[0;31mtest-host-004.example.com  Process exited with status 1 after 12s\033[0m\n",
		"test-host-005.example.com  Text on stdout without newline\n\x1b[0;31mtest-host-005.example.com  \x1b[0mText on stderr\n  More text\n",
	}
	for i, r := range results {
		if s := testformatter.formatOutput(r, 0); s != expected[i] {
			t.Errorf("Result %d, expected output %s, got %s", i, strconv.Quote(expected[i]), strconv.Quote(s))
		}
	}
}

func TestPrettyFormatterFormatResult(t *testing.T) {
	expected := []string{
		"\033[0;31mtest-host-001.example.com  It's always DNS after 12s\033[0m\n",
		"\033[0;32mtest-host-002.example.com  completed successfully after 12s\033[0m\n    May the forks be with you\n    And you\n",
		"\033[0;32mtest-host-003.example.com  completed successfully after 12s\033[0m\n    Newline is added automatically\n",
		"\033[0;31mtest-host-004.example.com  Process exited with status 1 after 12s\033[0m\n\033[0;90m----\033[0m\n    Text on stderr\n    More text\n",
		"\x1b[0;32mtest-host-005.example.com  completed successfully after 12s\x1b[0m\n    Text on stdout without newline\n\x1b[0;90m----\x1b[0m\n    Text on stderr\n    More text\n",
	}
	for i, r := range results {
		if s := testformatter.formatResult(r, 0); s != expected[i] {
			t.Errorf("Result %d, expected result %s, got %s", i, strconv.Quote(expected[i]), strconv.Quote(s))
		}
	}
}

func TestIndent(t *testing.T) {
	t.Skip("Not yet implemented")
}

func TestLogrusFormatter(t *testing.T) {
	t.Skip("Not yet implemented")
}
