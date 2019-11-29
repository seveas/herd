package scripting

import (
	"regexp"
	"testing"

	"github.com/go-test/deep"
	"github.com/seveas/katyusha"
)

func init() {
	deep.CompareUnexportedFields = true
}

func TestParseCommandLine(t *testing.T) {
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
	expected := [][]command{
		{
			addHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{}},
		},
		{},
		{
			addHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}},
		},
		{
			addHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}, {Name: "baz", Value: "quux", FuzzyTyping: true}}},
		},
		{
			addHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}},
			addHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "baz", Value: "quux", FuzzyTyping: true}}},
		},
		{
			addHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}},
			removeHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "baz", Value: "quux", FuzzyTyping: true}}},
		},
		{
			addHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}},
			removeHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "baz", Value: "quux", FuzzyTyping: true}}},
			addHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "zoinks", Value: "floop", FuzzyTyping: true}}},
		},
		{},
		{
			addHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true, Negate: true}}},
		},
		{
			addHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "foo", Value: regexp.MustCompile("bar"), Regex: true}}},
		},
		{
			addHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "foo", Value: regexp.MustCompile("bar"), Regex: true, Negate: true}}},
		},
		{
			addHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}},
		},
		{
			addHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "foo", Value: "bar", FuzzyTyping: true}}},
			addHostsCommand{glob: "*", attributes: katyusha.MatchAttributes{{Name: "baz", Value: "quux", FuzzyTyping: true}}},
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
		e := NewScriptEngine(nil, nil)
		err := e.ParseCommandLine(test, -1)
		if (err != nil && err.Error() != errors[i]) || (err == nil && errors[i] != "") {
			t.Errorf("Unexpected error %v, expected %v", err, errors[i])
		}
		if diff := deep.Equal(e.commands, expected[i]); diff != nil {
			t.Error(diff)
		}
		if errors[i] != "" {
			continue
		}
		test = append(test, "id", "seveas")
		e = NewScriptEngine(nil, nil)
		err = e.ParseCommandLine(test, len(test)-2)
		if (err != nil && err.Error() != errors[i]) || (err == nil && errors[i] != "") {
			t.Errorf("Unexpected error %v, expected %v", err, errors[i])
		}
		if diff := deep.Equal(e.commands, append(expected[i], runCommand{command: "id seveas"})); diff != nil {
			t.Error(diff)
		}
	}
}
