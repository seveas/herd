package katyusha

import (
	"fmt"
)

type MultiError struct {
	Errors  []error
	Message string
}

func (m *MultiError) Error() string {
	return m.Message
}

func (m *MultiError) Add(e error) {
	m.Errors = append(m.Errors, e)
	if m.Message == "" {
		m.Message = e.Error()
	} else {
		m.Message = fmt.Sprintf("%s\n%s", m.Message, e.Error())
	}
}

func (m *MultiError) AddHidden(e error) {
	m.Errors = append(m.Errors, e)
}
