package scripting

import (
	"regexp"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/seveas/herd"
)

func init() {
	deep.CompareUnexportedFields = true
}

func TestParseCommandLine(t *testing.T) {
	tests := []struct {
		spec []string
		cmd  []command
		err  string
	}{
		{
			spec: []string{"*"},
			cmd:  []command{addHostsCommand{glob: "*", attributes: herd.MatchAttributes{}, sampled: []string{}}},
			err:  "",
		},
		{
			spec: []string{"+", "*"},
			cmd:  []command{},
			err:  "incorrect filter: *",
		},
		{
			spec: []string{"*", "foo=bar"},
			cmd:  []command{addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}, sampled: []string{}}},
			err:  "",
		},
		{
			spec: []string{"*", "foo=bar", "baz=quux"},
			cmd:  []command{addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}, {Name: "baz", Value: "quux", FuzzyTyping: true}}, sampled: []string{}}},
			err:  "",
		},
		{
			spec: []string{"*", "foo=bar", "+", "*", "baz=quux"},
			cmd: []command{
				addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}, sampled: []string{}},
				addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "baz", Value: "quux", FuzzyTyping: true}}, sampled: []string{}},
			},
			err: "",
		},
		{
			spec: []string{"*", "foo=bar", "-", "*", "baz=quux"},
			cmd: []command{
				addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}, sampled: []string{}},
				removeHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "baz", Value: "quux", FuzzyTyping: true}}},
			},
			err: "",
		},
		{
			spec: []string{"*", "foo=bar", "-", "*", "baz=quux", "+", "*", "zoinks=floop"},
			cmd: []command{
				addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}, sampled: []string{}},
				removeHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "baz", Value: "quux", FuzzyTyping: true}}},
				addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "zoinks", Value: "floop", FuzzyTyping: true}}, sampled: []string{}},
			},
			err: "",
		},
		{
			spec: []string{"*", "foo"},
			cmd:  []command{},
			err:  "incorrect filter: foo",
		},
		{
			spec: []string{"*", "foo!=bar"},
			cmd:  []command{addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true, Negate: true}}, sampled: []string{}}},
			err:  "",
		},
		{
			spec: []string{"*", "foo=~bar"},
			cmd:  []command{addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: regexp.MustCompile("bar"), Regex: true}}, sampled: []string{}}},
			err:  "",
		},
		{
			spec: []string{"*", "foo!~bar"},
			cmd:  []command{addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: regexp.MustCompile("bar"), Regex: true, Negate: true}}, sampled: []string{}}},
			err:  "",
		},
		{
			spec: []string{"foo=bar"},
			cmd:  []command{addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}, sampled: []string{}}},
			err:  "",
		},
		{
			spec: []string{"foo=bar", "+", "baz=quux"},
			cmd: []command{
				addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}, sampled: []string{}},
				addHostsCommand{glob: "*", attributes: herd.MatchAttributes{{Name: "baz", Value: "quux", FuzzyTyping: true}}, sampled: []string{}},
			},
			err: "",
		},
		{
			spec: []string{"foo:1"},
			cmd:  []command{addHostsCommand{glob: "*", attributes: herd.MatchAttributes{}, sampled: []string{"foo"}, count: 1}},
			err:  "",
		},
		{
			spec: []string{"foo:bar:1"},
			cmd:  []command{addHostsCommand{glob: "*", attributes: herd.MatchAttributes{}, sampled: []string{"foo", "bar"}, count: 1}},
			err:  "",
		},
		{
			spec: []string{":::foo:bar::baz:::quux:3"},
			cmd:  []command{addHostsCommand{glob: "*", attributes: herd.MatchAttributes{}, sampled: []string{"foo", "bar:baz:", "quux"}, count: 3}},
			err:  "",
		},
		{
			spec: []string{"foo:1", "bar:1"},
			cmd:  []command{},
			err:  "only one sampling per hostspec allowed",
		},
	}

	for _, test := range tests {
		t.Run(strings.Join(test.spec, " "), func(t *testing.T) {
			e := NewScriptEngine(nil, nil)
			err := e.ParseCommandLine(test.spec, -1)
			if (err != nil && err.Error() != test.err) || (err == nil && test.err != "") {
				t.Errorf("Unexpected error %v, expected %v", err, test.err)
			}
			if diff := deep.Equal(test.cmd, e.commands); diff != nil {
				t.Error(diff)
			}
			if test.err != "" {
				return
			}
			test.spec = append(test.spec, "id", "seveas")
			e = NewScriptEngine(nil, nil)
			err = e.ParseCommandLine(test.spec, len(test.spec)-2)
			if (err != nil && err.Error() != test.err) || (err == nil && test.err != "") {
				t.Errorf("Unexpected error %v, expected %v", err, test.err)
			}
			if diff := deep.Equal(e.commands, append(test.cmd, runCommand{command: "id seveas"})); diff != nil {
				t.Error(diff)
			}
		})
	}
}
