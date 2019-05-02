package cmd

import (
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
	}
	expected := [][]herd.Command{
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.HostAttributes{}},
		},
		nil,
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.HostAttributes{"foo": "bar"}},
		},
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.HostAttributes{"foo": "bar", "baz": "quux"}},
		},
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.HostAttributes{"foo": "bar"}},
			herd.AddHostsCommand{Glob: "*", Attributes: herd.HostAttributes{"baz": "quux"}},
		},
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.HostAttributes{"foo": "bar"}},
			herd.RemoveHostsCommand{Glob: "*", Attributes: herd.HostAttributes{"baz": "quux"}},
		},
		{
			herd.AddHostsCommand{Glob: "*", Attributes: herd.HostAttributes{"foo": "bar"}},
			herd.RemoveHostsCommand{Glob: "*", Attributes: herd.HostAttributes{"baz": "quux"}},
			herd.AddHostsCommand{Glob: "*", Attributes: herd.HostAttributes{"zoinks": "floop"}},
		},
		nil,
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
