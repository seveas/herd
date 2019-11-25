package scripting

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/go-test/deep"
	"github.com/seveas/herd"
)

type testcase struct {
	program  string
	commands []Command
	errors   []error
	err      error
}

var testcases = []testcase{
	{
		program:  "",
		commands: []Command{},
	},
	{
		program: "syntax error",
		errors:  []error{fmt.Errorf("line 1:0 mismatched input 'syntax' expecting {<EOF>, <NEWLINE>, RUN, 'set', 'add', 'remove', 'list'}")},
	},
	{
		program: "add hosts * foo == \"bar\"\n" +
			"remove hosts * foo == \"bar\"\n" +
			"add hosts * foo == 1\n" +
			"add hosts * foo == nil\n" +
			"add hosts * foo =~ /bar/\n" +
			"list hosts\n" +
			"# comment, should be ignored\n" +
			"run find / -name 'whatever' -delete\n" +
			"list hosts --oneline\n",
		commands: []Command{
			AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: "bar"}}},
			RemoveHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: "bar"}}},
			AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: int64(1)}}},
			AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: nil}}},
			AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: regexp.MustCompile("bar"), Regex: true}}},
			ListHostsCommand{},
			RunCommand{Command: "find / -name 'whatever' -delete"},
			ListHostsCommand{OneLine: true},
		},
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
		commands, err := ParseCode(tc.program)
		if diff := deep.Equal(tc.commands, commands); diff != nil {
			t.Errorf("(%d) Unexpected diff in commands:\n%s", i, diff)
		}
		if diff := deep.Equal(tc.err, err); diff != nil {
			t.Errorf("(%d) Unexpected diff in error:\n%s", i, diff)
		}
	}
}
