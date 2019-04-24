package katyusha

import (
	"fmt"
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

type RemoveHostsCommand struct {
	Glob       string
	Attributes HostAttributes
}

func (c RemoveHostsCommand) Execute(r *Runner) error {
	r.RemoveHosts(c.Glob, c.Attributes)
	return nil
}

func (c RemoveHostsCommand) String() string {
	return fmt.Sprintf("remove hosts %s %v", c.Glob, c.Attributes)
}

type ListHostsCommand struct {
	OneLine bool
}

func (c ListHostsCommand) Execute(r *Runner) error {
	r.ListHosts(c.OneLine)
	return nil
}

func (c ListHostsCommand) String() string {
	if c.OneLine {
		return "list hosts --oneline"
	} else {
		return "list hosts"
	}
}

type RunCommand struct {
	Command string
}

func (c RunCommand) Execute(r *Runner) error {
	hi := r.Run(c.Command)
	UI.PrintHistoryItem(hi)
	return nil
}

func (c RunCommand) String() string {
	return "run " + c.Command
}
