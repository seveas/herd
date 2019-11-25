package scripting

import (
	"fmt"
	"strings"

	"github.com/seveas/katyusha"
	"github.com/spf13/viper"
)

type Command interface {
	Execute(r *katyusha.Runner) error
}

type SetCommand struct {
	Variable string
	Value    interface{}
}

func (c SetCommand) Execute(r *katyusha.Runner) error {
	viper.Set(c.Variable, c.Value)
	return nil
}

func (c SetCommand) String() string {
	return fmt.Sprintf("set %s %v", c.Variable, c.Value)
}

type AddHostsCommand struct {
	Glob       string
	Attributes katyusha.MatchAttributes
}

func (c AddHostsCommand) Execute(r *katyusha.Runner) error {
	r.AddHosts(c.Glob, c.Attributes)
	return nil
}

func (c AddHostsCommand) String() string {
	return fmt.Sprintf("add hosts %s %v", c.Glob, c.Attributes)
}

type RemoveHostsCommand struct {
	Glob       string
	Attributes katyusha.MatchAttributes
}

func (c RemoveHostsCommand) Execute(r *katyusha.Runner) error {
	r.RemoveHosts(c.Glob, c.Attributes)
	return nil
}

func (c RemoveHostsCommand) String() string {
	return fmt.Sprintf("remove hosts %s %v", c.Glob, c.Attributes)
}

type ListHostsCommand struct {
	OneLine       bool
	AllAttributes bool
	Attributes    []string
	Csv           bool
}

func (c ListHostsCommand) Execute(r *katyusha.Runner) error {
	r.ListHosts(c.OneLine, c.AllAttributes, c.Attributes, c.Csv)
	return nil
}

func (c ListHostsCommand) String() string {
	ret := "list hosts"
	if c.OneLine {
		ret += " --oneline"
	}
	if c.AllAttributes {
		ret += " --all-attributes"
	}
	if len(c.Attributes) != 0 {
		ret += " --attributes=" + strings.Join(c.Attributes, ",")
	}
	return ret
}

type RunCommand struct {
	Command string
}

func (c RunCommand) Execute(r *katyusha.Runner) error {
	r.Run(c.Command)
	return nil
}

func (c RunCommand) String() string {
	return "run " + c.Command
}
