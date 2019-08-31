package herd

import (
	"fmt"

	"github.com/spf13/viper"
)

type Command interface {
	Execute(r *Runner) error
}

type FilterCommand interface {
	Match(h *Host) bool
}

type SetCommand struct {
	Variable string
	Value    interface{}
}

func (c SetCommand) Execute(r *Runner) error {
	viper.Set(c.Variable, c.Value)
	return nil
}

func (c SetCommand) String() string {
	return fmt.Sprintf("set %s %v", c.Variable, c.Value)
}

type AddHostsCommand struct {
	Glob       string
	Attributes MatchAttributes
}

func (c AddHostsCommand) Execute(r *Runner) error {
	r.AddHosts(c.Glob, c.Attributes)
	return nil
}

func (c AddHostsCommand) Match(h *Host) bool {
	return h.Match(c.Glob, c.Attributes)
}

func (c AddHostsCommand) String() string {
	return fmt.Sprintf("add hosts %s %v", c.Glob, c.Attributes)
}

type RemoveHostsCommand struct {
	Glob       string
	Attributes MatchAttributes
}

func (c RemoveHostsCommand) Execute(r *Runner) error {
	r.RemoveHosts(c.Glob, c.Attributes)
	return nil
}

func (c RemoveHostsCommand) Match(h *Host) bool {
	return !h.Match(c.Glob, c.Attributes)
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
	r.Run(c.Command)
	return nil
}

func (c RunCommand) String() string {
	return "run " + c.Command
}
