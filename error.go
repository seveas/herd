package katyusha

import (
	"fmt"
	"strings"
)

type MultiError struct {
	Subject  string
	errors   []error
	messages []string
}

func (m *MultiError) Error() string {
	if m.Subject != "" && len(m.messages) != 0 {
		return fmt.Sprintf("%s:\n%s", m.Subject, strings.Join(m.messages, "\n"))
	}
	return strings.Join(m.messages, "\n")
}

func (m *MultiError) Add(e error) {
	m.errors = append(m.errors, e)
	m.messages = append(m.messages, e.Error())
}

func (m *MultiError) HasErrors() bool {
	return len(m.errors) > 0
}
