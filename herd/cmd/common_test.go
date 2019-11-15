package cmd

import (
	"regexp"
	"testing"

	"github.com/go-test/deep"
	"github.com/seveas/herd"
	"github.com/spf13/cobra"
)

func TestSplitArgs(t *testing.T) {
	cmd := &cobra.Command{}

	tests := [][]string{
		{"arg1", "arg2", "--", "arg3", "arg4"},
		{"arg1", "arg2", "--"},
		{"--", "arg3", "arg4"},
	}
	expected := [][][]string{
		{{"arg1", "arg2"}, {"arg3", "arg4"}},
		{{"arg1", "arg2"}, {}},
		{{}, {"arg3", "arg4"}},
	}

	for i, test := range tests {
		cmd.ParseFlags(test)
		before, after := splitArgs(cmd, cmd.Flags().Args())
		if diff := deep.Equal([][]string{before, after}, expected[i]); diff != nil {
			t.Error(diff)
		}
	}
}

func TestFilterCommands(t *testing.T) {
	tests := [][]string{
		{"*"},
		{"+", "*"},
		{"*", "foo=bar"},
		{"*", "foo=bar", "baz=quux"},
		{"*", "foo=bar", "+", "*", "baz=quux"},
		{"*", "foo=bar", "-", "*", "baz=quux"},
		{"*", "foo=bar", "-", "*", "baz=quux", "+", "*", "zoinks=floop"},
		{"*", "foo"},
		{"*", "foo!=bar"},
		{"*", "foo=~bar"},
		{"*", "foo!~bar"},
		{"foo=bar"},
		{"foo=bar", "+", "baz=quux"},
	}
	expected := [][]herd.Command{
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{}},
		},
		nil,
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}},
		},
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}, {Name: "baz", Value: "quux", FuzzyTyping: true}}},
		},
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}},
			herd.AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "baz", Value: "quux", FuzzyTyping: true}}},
		},
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}},
			herd.RemoveHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "baz", Value: "quux", FuzzyTyping: true}}},
		},
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}},
			herd.RemoveHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "baz", Value: "quux", FuzzyTyping: true}}},
			herd.AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "zoinks", Value: "floop", FuzzyTyping: true}}},
		},
		nil,
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true, Negate: true}}},
		},
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: regexp.MustCompile("bar"), Regex: true}}},
		},
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: regexp.MustCompile("bar"), Regex: true, Negate: true}}},
		},
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}},
		},
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}},
			herd.AddHostsCommand{Glob: "*", Attributes: herd.MatchAttributes{{Name: "baz", Value: "quux", FuzzyTyping: true}}},
		},
	}
	errors := []string{
		"",
		"incorrect filter: *",
		"",
		"",
		"",
		"",
		"",
		"incorrect filter: foo",
		"",
		"",
		"",
		"",
		"",
	}

	for i, test := range tests {
		commands, err := filterCommands(test)
		if (err != nil && err.Error() != errors[i]) || (err == nil && errors[i] != "") {
			t.Errorf("Unexpected error %v, expected %v", err, errors[i])
		}
		if diff := deep.Equal(commands, expected[i]); diff != nil {
			t.Error(diff)
		}
	}
}
