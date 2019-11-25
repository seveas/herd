package katyusha

import (
	"fmt"
	"testing"
)

func TestEmptyError(t *testing.T) {
	e := &MultiError{Subject: "Test error"}
	if e.HasErrors() {
		t.Error("Empty multierror seems to have an error")
	}
	if e.Error() != "" {
		t.Errorf("Empty multierror has a non-empty string: %s", e.Error())
	}
}

func TestNoSubject(t *testing.T) {
	e := &MultiError{}
	e.Add(fmt.Errorf("Test error 1"))
	if e.Error() != "Test error 1" {
		t.Errorf("Unexpected single-error string: %s", e.Error())
	}
	e.Add(fmt.Errorf("Test error 2"))
	if e.Error() != "Test error 1\nTest error 2" {
		t.Errorf("Unexpected single-error string: %s", e.Error())
	}
}

func TestWithSubject(t *testing.T) {
	e := &MultiError{Subject: "Test subject"}
	e.Add(fmt.Errorf("Test error 1"))
	if e.Error() != "Test subject:\nTest error 1" {
		t.Errorf("Unexpected single-error string: %s", e.Error())
	}
	e.Add(fmt.Errorf("Test error 2"))
	if e.Error() != "Test subject:\nTest error 1\nTest error 2" {
		t.Errorf("Unexpected single-error string: %s", e.Error())
	}
}
