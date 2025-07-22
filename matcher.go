package herd

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

type MatchAttribute struct {
	Name        string
	FuzzyTyping bool
	Negate      bool
	Regex       bool
	Value       interface{}
}

func (m MatchAttribute) String() string {
	c1, c2 := '=', '='
	f := ""
	if m.FuzzyTyping {
		f = " (≈)"
	}
	if m.Negate {
		c1 = '!'
	}
	if m.Regex {
		c2 = '~'
	}
	return fmt.Sprintf("%v %c%c %v%s", m.Name, c1, c2, m.Value, f)
}

func (m MatchAttribute) Match(value interface{}) (matches bool) {
	defer func() {
		if m.Negate {
			matches = !matches
		}
	}()
	if svalue := reflect.ValueOf(value); svalue.Kind() == reflect.Slice {
		// Here we ignore Negate to make sure we filter for any/none matching
		mx := MatchAttribute{Name: m.Name, Value: m.Value, FuzzyTyping: m.FuzzyTyping, Regex: m.Regex}
		for i := 0; i < svalue.Len(); i++ {
			if mx.Match(svalue.Index(i).Interface()) {
				return true
			}
		}
		return false
	}
	if m.Value == value {
		return true
	}
	if m.Regex {
		svalue, ok := value.(string)
		return ok && m.Value.(*regexp.Regexp).MatchString(svalue)
	}
	if m.FuzzyTyping {
		if bvalue, ok := value.(bool); ok && (m.Value == "true" || m.Value == "false") {
			return bvalue == (m.Value == "true")
		}
		if m.Value == "nil" {
			return value == nil
		}
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			myival, err := strconv.ParseInt(m.Value.(string), 0, 64)
			if err != nil {
				return false
			}
			return v.Int() == myival
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			myival, err := strconv.ParseUint(m.Value.(string), 0, 64)
			if err != nil {
				return false
			}
			return v.Uint() == myival
		}
	}
	// Let's be gentle on all the int types in attributes
	if myival, ok := m.Value.(int64); ok {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return v.Int() == myival
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int64(v.Uint()) == myival // nolint:gosec // FIXME I need to take another look at this. Security implications are limited as this only affects which hosts to match.
		}
	}
	return false
}

type MatchAttributes []MatchAttribute
