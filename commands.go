package katyusha

import (
	"fmt"
	"os"
	"reflect"
)

type Command interface {
	Execute(r *Runner) error
}

type SetCommand struct {
	VariableName string
	// Pointer to the actual variable
	Variable interface{}
	// New value. The setter code *must* make sure that it's set to something of the correct type
	Value interface{}
}

func (c SetCommand) Execute(r *Runner) error {
	reflect.ValueOf(c.Variable).Elem().Set(reflect.ValueOf(c.Value))
	return nil
}

func (c SetCommand) String() string {
	return fmt.Sprintf("set %s %v", c.VariableName, c.Value)
}

type AddHostsCommand struct {
	Glob       string
	Attributes HostAttributes
}

func (c AddHostsCommand) Execute(r *Runner) error {
	r.AddHosts(c.Glob, c.Attributes)
	return nil
}

func (c AddHostsCommand) String() string {
	return fmt.Sprintf("add hosts %s %v", c.Glob, c.Attributes)
}

type RunCommand struct {
	Command   string
	Formatter Formatter
}

func (c RunCommand) Execute(r *Runner) error {
	hi := r.Run(c.Command)
	// FIXME should go through the UI layer
	c.Formatter.Format(hi, os.Stdout)
	return nil
}

func (c RunCommand) String() string {
	return "run " + c.Command
}
