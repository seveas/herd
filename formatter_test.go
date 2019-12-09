package herd

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var results []*Result

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
	if x := Formatters["pretty"].FormatCommand("hello world"); x != "\033[36mhello world\033[0m\n" {
		t.Errorf("Expected colored string, got %s", strconv.Quote(x))
	}
}

func TestPrettyFormatterFormatStatus(t *testing.T) {
	f := Formatters["pretty"]
	expected := []string{
		"\033[31mtest-host-001.example.com  It's always DNS after 12s\033[0m\n",
		"\033[32mtest-host-002.example.com  completed successfully after 12s\033[0m\n",
		"\033[32mtest-host-003.example.com  completed successfully after 12s\033[0m\n",
		"\033[31mtest-host-004.example.com  Process exited with status 1 after 12s\033[0m\n",
		"\033[32mtest-host-005.example.com  completed successfully after 12s\033[0m\n",
	}
	for i, r := range results {
		if s := f.FormatStatus(r, 0); s != expected[i] {
			t.Errorf("Result %d, expected status %s, got %s", i, strconv.Quote(expected[i]), strconv.Quote(s))
		}
	}
}

func TestPrettyFormatterFormatOutput(t *testing.T) {
	f := Formatters["pretty"]
	expected := []string{
		"\033[31mtest-host-001.example.com  It's always DNS after 12s\033[0m\n",
		"test-host-002.example.com  May the forks be with you\n  And you\n",
		"test-host-003.example.com  Newline is added automatically\n",
		"\033[31mtest-host-004.example.com  \033[0mText on stderr\n  More text\n\033[31mtest-host-004.example.com  Process exited with status 1 after 12s\033[0m\n",
		"test-host-005.example.com  Text on stdout without newline\n\x1b[31mtest-host-005.example.com  \x1b[0mText on stderr\n  More text\n",
	}
	for i, r := range results {
		if s := f.FormatOutput(r, 0); s != expected[i] {
			t.Errorf("Result %d, expected output %s, got %s", i, strconv.Quote(expected[i]), strconv.Quote(s))
		}
	}
}

func TestPrettyFormatterFormatResult(t *testing.T) {
	f := Formatters["pretty"]
	expected := []string{
		"\033[31mtest-host-001.example.com  It's always DNS after 12s\033[0m\n",
		"\033[32mtest-host-002.example.com  completed successfully after 12s\033[0m\n    May the forks be with you\n    And you\n",
		"\033[32mtest-host-003.example.com  completed successfully after 12s\033[0m\n    Newline is added automatically\n",
		"\033[31mtest-host-004.example.com  Process exited with status 1 after 12s\033[0m\n\033[90m----\033[0m\n    Text on stderr\n    More text\n",
		"\x1b[32mtest-host-005.example.com  completed successfully after 12s\x1b[0m\n    Text on stdout without newline\n\x1b[90m----\x1b[0m\n    Text on stderr\n    More text\n",
	}
	for i, r := range results {
		if s := f.FormatResult(r); s != expected[i] {
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
