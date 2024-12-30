package herd

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var results []*Result

var testformatter = newPrettyFormatter(defaultColorConfigDark)

func init() {
	i := 0
	count := func() int {
		i++
		return i
	}
	start := time.Date(2019, 12, 8, 20, 26, 0, 0, time.UTC)
	results = []*Result{
		{
			Host:        fmt.Sprintf("test-host-%03d.example.com", count()),
			Err:         fmt.Errorf("It's always DNS"),
			ExitStatus:  -1,
			Stdout:      []byte{},
			Stderr:      []byte{},
			StartTime:   start,
			EndTime:     start.Add(12 * time.Second),
			ElapsedTime: 12,
		},
		{
			Host:        fmt.Sprintf("test-host-%03d.example.com", count()),
			Stdout:      []byte("May the forks be with you\nAnd you\n"),
			Stderr:      []byte{},
			StartTime:   start,
			EndTime:     start.Add(12 * time.Second),
			ElapsedTime: 12,
		},
		{
			Host:        fmt.Sprintf("test-host-%03d.example.com", count()),
			Stdout:      []byte("Newline is added automatically"),
			Stderr:      []byte{},
			StartTime:   start,
			EndTime:     start.Add(12 * time.Second),
			ElapsedTime: 12,
		},
		{
			Host:        fmt.Sprintf("test-host-%03d.example.com", count()),
			Err:         fmt.Errorf("exited with status 1"),
			ExitStatus:  1,
			Stdout:      []byte{},
			Stderr:      []byte("Text on stderr\nMore text\n"),
			StartTime:   start,
			EndTime:     start.Add(12 * time.Second),
			ElapsedTime: 12,
		},
		{
			Host:        fmt.Sprintf("test-host-%03d.example.com", count()),
			Err:         fmt.Errorf("exited with status 1"),
			ExitStatus:  1,
			Stdout:      []byte("Text on stdout\n"),
			Stderr:      []byte("Text on stderr\nMore text\n"),
			StartTime:   start,
			EndTime:     start.Add(12 * time.Second),
			ElapsedTime: 12,
		},
		{
			Host:        fmt.Sprintf("test-host-%03d.example.com", count()),
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
		"\033[0;33mtest-host-004.example.com  exited with status 1 after 12s\033[0m\n",
		"\033[0;33mtest-host-005.example.com  exited with status 1 after 12s\033[0m\n",
		"\033[0;32mtest-host-006.example.com  completed successfully after 12s\033[0m\n",
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
		"\033[0;32mtest-host-002.example.com  \033[0mMay the forks be with you\n  And you\n",
		"\033[0;32mtest-host-003.example.com  \033[0mNewline is added automatically\n",
		"\033[0;33mtest-host-004.example.com  \033[0mText on stderr\n  More text\n\033[0;33mtest-host-004.example.com  exited with status 1 after 12s\033[0m\n",
		"\033[0;33mtest-host-005.example.com  \033[0mText on stdout\n\033[0;33mtest-host-005.example.com  \033[0mText on stderr\n  More text\n\033[0;33mtest-host-005.example.com  exited with status 1 after 12s\033[0m\n",
		"\033[0;32mtest-host-006.example.com  \033[0mText on stdout without newline\n\x1b[0;33mtest-host-006.example.com  \x1b[0mText on stderr\n  More text\n",
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
		"\033[0;33mtest-host-004.example.com  exited with status 1 after 12s\033[0m\n\033[0;90m----\033[0m\n    Text on stderr\n    More text\n",
		"\033[0;33mtest-host-005.example.com  exited with status 1 after 12s\033[0m\n    Text on stdout\n\033[0;90m----\033[0m\n    Text on stderr\n    More text\n",
		"\033[0;32mtest-host-006.example.com  completed successfully after 12s\x1b[0m\n    Text on stdout without newline\n\x1b[0;90m----\x1b[0m\n    Text on stderr\n    More text\n",
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
