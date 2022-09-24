package herd

import (
	"regexp"
	"testing"
)

type testcase struct {
	a MatchAttribute
	v interface{}
	m bool
}

func TestMatchAttribute(t *testing.T) {
	testcases := []testcase{
		// Basic equality
		{a: MatchAttribute{Value: 1}, v: 1, m: true},
		{a: MatchAttribute{Value: 1}, v: 0, m: false},
		{a: MatchAttribute{Value: "1"}, v: "1", m: true},
		{a: MatchAttribute{Value: "1"}, v: "0", m: false},
		{a: MatchAttribute{Value: true}, v: true, m: true},
		{a: MatchAttribute{Value: true}, v: false, m: false},
		{a: MatchAttribute{Value: false}, v: nil, m: false},
		// Unequal types
		{a: MatchAttribute{Value: "1"}, v: 1, m: false},
		{a: MatchAttribute{Value: "true"}, v: true, m: false},
		{a: MatchAttribute{Value: false}, v: nil, m: false},
		// Integer niceness
		{a: MatchAttribute{Value: int64(1)}, v: int(1), m: true},
		// String fuzziness
		{a: MatchAttribute{Value: "1", FuzzyTyping: true}, v: 1, m: true},
		{a: MatchAttribute{Value: "1", FuzzyTyping: true}, v: 0, m: false},
		{a: MatchAttribute{Value: "1", FuzzyTyping: true}, v: int32(1), m: true},
		{a: MatchAttribute{Value: "1", FuzzyTyping: true}, v: uint16(1), m: true},
		{a: MatchAttribute{Value: "true", FuzzyTyping: true}, v: true, m: true},
		{a: MatchAttribute{Value: "nil", FuzzyTyping: true}, v: nil, m: true},
		// Regular expressions
		{a: MatchAttribute{Value: regexp.MustCompile("hello"), Regex: true}, v: "hello", m: true},
		{a: MatchAttribute{Value: regexp.MustCompile("hello"), Regex: true}, v: "hello world", m: true},
		{a: MatchAttribute{Value: regexp.MustCompile("hello$"), Regex: true}, v: "hello world", m: false},
		// Slices
		{a: MatchAttribute{Value: 1}, v: []int{2, 3, 1, 4}, m: true},
		{a: MatchAttribute{Value: 1}, v: []int{2, 3, 4}, m: false},
		{a: MatchAttribute{Value: 1}, v: []int{1}, m: true},
		{a: MatchAttribute{Value: 1}, v: []int{}, m: false},
	}

	for i, c := range testcases {
		c.a.Name = "v"
		if m := c.a.Match(c.v); m != c.m {
			if c.m {
				t.Errorf("(%d) expected %v (%T) to match %v (%T), but they did not match", i, c.a, c.a.Value, c.v, c.v)
			} else {
				t.Errorf("(%d) expected %v (%T) to not match %v (%T), but they did match", i, c.a, c.a.Value, c.v, c.v)
			}
		}
		// Test the negation as well
		a := MatchAttribute{Name: "v", FuzzyTyping: c.a.FuzzyTyping, Value: c.a.Value, Regex: c.a.Regex, Negate: !c.a.Negate}
		if m := a.Match(c.v); m != !c.m {
			if !c.m {
				t.Errorf("(%d) expected %v (%T) to match %v (%T), but they did not match", i, a, a.Value, c.v, c.v)
			} else {
				t.Errorf("(%d) expected %v (%T) to not match %v (%T), but they did match", i, a, a.Value, c.v, c.v)
			}
		}
	}
}
