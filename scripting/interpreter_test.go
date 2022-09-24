package scripting

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/seveas/herd"
	"github.com/sirupsen/logrus"
)

func init() {
	deep.CompareUnexportedFields = true
}

type testcase struct {
	program  string
	commands []command
	errors   []error
	err      error
}

var testcases = []testcase{
	{
		program:  "",
		commands: []command{},
	},
	{
		program: "syntax error",
		errors:  []error{fmt.Errorf("line 1:0 mismatched input 'syntax' expecting {<EOF>, '\\n', RUN, 'set', 'add', 'remove', 'list'}")},
	},
	{
		program: strings.Join([]string{
			"add hosts * foo == \"bar\"",
			"remove hosts * foo == \"bar\"",
			"add hosts * foo == 1",
			"add hosts * foo == nil",
			"add hosts * foo =~ /bar/",
			"list hosts",
			"# comment, should be ignored",
			"run find / -name 'whatever' -delete",
			"list hosts {OneLine: true, Align: false, Separator: \"-\", AllAttributes: true, Attributes: [\"foo\"], Csv: true, Header: false}",
		}, "\n") + "\n",
		commands: []command{
			addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: "bar"}}},
			removeHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: "bar"}}},
			addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: int64(1)}}},
			addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: nil}}},
			addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: regexp.MustCompile("bar"), Regex: true}}},
			listHostsCommand{opts: herd.HostListOptions{Separator: ",", Header: true, Align: true}},
			runCommand{command: "find / -name 'whatever' -delete"},
			listHostsCommand{opts: herd.HostListOptions{Separator: "-", OneLine: true, AllAttributes: true, Attributes: []string{"foo"}, Csv: true}},
		},
	},
	{
		program: "list hosts {OneLine: 1, Align: \"foo\", Separator: nil, AllAttributes: 3m, Attributes: [true], Csv: 21, Header: \"oink\"}\n",
		errors: []error{
			fmt.Errorf("line 1:0 OneLine must be a boolean value"),
			fmt.Errorf("line 1:0 Csv must be a boolean value"),
			fmt.Errorf("line 1:0 Align must be a boolean value"),
			fmt.Errorf("line 1:0 AllAttributes must be a boolean value"),
			fmt.Errorf("line 1:0 Header must be a boolean value"),
			fmt.Errorf("line 1:0 Separator must be a string"),
			fmt.Errorf("line 1:0 Attributes must be a list of strings"),
		},
	},
	{
		program: strings.Join([]string{
			"set Splay 5s",
			"set Timeout 1m",
			"set HostTimeout 10m",
			"set ConnectTimeout 10s",
			"set Parallel 50",
			"set Timestamp true",
			"set NoPager true",
			"set NoColor true",
			"set LogLevel \"debug\"",
			"set Output \"inline\"",
		}, "\n") + "\n",
		commands: []command{
			setCommand{variable: "Splay", value: 5 * time.Second},
			setCommand{variable: "Timeout", value: 1 * time.Minute},
			setCommand{variable: "HostTimeout", value: 10 * time.Minute},
			setCommand{variable: "ConnectTimeout", value: 10 * time.Second},
			setCommand{variable: "Parallel", value: int64(50)},
			setCommand{variable: "Timestamp", value: true},
			setCommand{variable: "NoPager", value: true},
			setCommand{variable: "NoColor", value: true},
			setCommand{variable: "LogLevel", value: logrus.DebugLevel},
			setCommand{variable: "Output", value: herd.OutputInline},
		},
	},
	{
		program: "set Splay true\n",
		errors:  []error{fmt.Errorf("line 1:10 Splay must be a duration")},
	},
	{
		program: "set Timeout true\n",
		errors:  []error{fmt.Errorf("line 1:12 Timeout must be a duration")},
	},
	{
		program: "set HostTimeout true\n",
		errors:  []error{fmt.Errorf("line 1:16 HostTimeout must be a duration")},
	},
	{
		program: "set ConnectTimeout false\n",
		errors:  []error{fmt.Errorf("line 1:19 ConnectTimeout must be a duration")},
	},
	{
		program: "set Parallel \"nope\"\n",
		errors:  []error{fmt.Errorf("line 1:13 Parallel must be a number")},
	},
	{
		program: "set Timestamp 23\n",
		errors:  []error{fmt.Errorf("line 1:14 Timestamp must be a boolean")},
	},
	{
		program: "set NoPager 42\n",
		errors:  []error{fmt.Errorf("line 1:12 NoPager must be a boolean")},
	},
	{
		program: "set NoColor 71\n",
		errors:  []error{fmt.Errorf("line 1:12 NoColor must be a boolean")},
	},
	{
		program: "set Output false\n",
		errors:  []error{fmt.Errorf("line 1:11 Output must be a string")},
	},
	{
		program: "set Output \"foo\"\n",
		errors:  []error{fmt.Errorf("line 1:11 Unknown output mode: foo. Known modes: all, per-host, inline, tail")},
	},
	{
		program: "set LogLevel nil\n",
		errors:  []error{fmt.Errorf("line 1:13 LogLevel must be a string")},
	},
	{
		program: "set LogLevel \"foo\"\n",
		errors:  []error{fmt.Errorf("line 1:13 Unknown loglevel: foo. Known loglevels: DEBUG, INFO, NORMAL, WARNING, ERROR")},
	},
	{
		program: "set NonExistent true\n",
		errors:  []error{fmt.Errorf("line 1:4 Unknown variable: NonExistent")},
	},
	{
		program: "set timeout 10s\n",
		errors:  []error{fmt.Errorf("line 1:4 Unknown variable: timeout")},
	},
}

func TestScripts(t *testing.T) {
	for i, tc := range testcases {
		if tc.errors != nil {
			err := &herd.MultiError{Subject: "Syntax errors found"}
			for _, e := range tc.errors {
				err.Add(e)
			}
			tc.err = err
		}
		commands, err := parseCode(tc.program)
		if diff := deep.Equal(tc.commands, commands); diff != nil {
			t.Errorf("(%d) Unexpected diff in commands:\n%s", i, diff)
		}
		if diff := deep.Equal(tc.err, err); diff != nil {
			t.Errorf("(%d) Unexpected diff in error:\n%s", i, diff)
		}
	}
}
