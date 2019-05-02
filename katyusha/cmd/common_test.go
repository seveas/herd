package cmd

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/seveas/katyusha"
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
	expected := [][]katyusha.Command{
		{
			katyusha.AddHostsCommand{Glob: "*", Attributes: katyusha.HostAttributes{}},
		},
		nil,
		{
			katyusha.AddHostsCommand{Glob: "*", Attributes: katyusha.HostAttributes{"foo": "bar"}},
		},
		{
			katyusha.AddHostsCommand{Glob: "*", Attributes: katyusha.HostAttributes{"foo": "bar", "baz": "quux"}},
		},
		{
			katyusha.AddHostsCommand{Glob: "*", Attributes: katyusha.HostAttributes{"foo": "bar"}},
			katyusha.AddHostsCommand{Glob: "*", Attributes: katyusha.HostAttributes{"baz": "quux"}},
		},
		{
			katyusha.AddHostsCommand{Glob: "*", Attributes: katyusha.HostAttributes{"foo": "bar"}},
			katyusha.RemoveHostsCommand{Glob: "*", Attributes: katyusha.HostAttributes{"baz": "quux"}},
		},
		{
			katyusha.AddHostsCommand{Glob: "*", Attributes: katyusha.HostAttributes{"foo": "bar"}},
			katyusha.RemoveHostsCommand{Glob: "*", Attributes: katyusha.HostAttributes{"baz": "quux"}},
			katyusha.AddHostsCommand{Glob: "*", Attributes: katyusha.HostAttributes{"zoinks": "floop"}},
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
