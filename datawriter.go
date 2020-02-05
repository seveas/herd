package katyusha

import (
	"fmt"
	"io"
	"strings"
)

type datawriter interface {
	Write([]string) error
	Flush()
}

type columnizer struct {
	rows    [][]string
	lengths []int
	output  io.Writer
	sep     string
}

func newColumnizer(w io.Writer, sep string) *columnizer {
	return &columnizer{rows: make([][]string, 0), output: w, sep: sep}
}

func (c *columnizer) Write(r []string) error {
	if c.lengths == nil {
		c.lengths = make([]int, len(r))
	}
	c.rows = append(c.rows, r)
	for i, v := range r {
		if l := len(v); l > c.lengths[i] {
			c.lengths[i] = l
		}
	}
	return nil
}

func (c *columnizer) Flush() {
	for _, r := range c.rows {
		for i, v := range r {
			if i > 0 {
				fmt.Fprint(c.output, c.sep)
			}
			fmt.Fprintf(c.output, "%-*s", c.lengths[i], v)
		}
		fmt.Fprint(c.output, "\n")
	}
}

type passthrough struct {
	output io.Writer
}

func newPassthrough(w io.Writer) *passthrough {
	return &passthrough{output: w}
}

func (p *passthrough) Write(r []string) error {
	_, err := p.output.Write([]byte(strings.Join(r, " ") + "\n"))
	return err
}

func (p *passthrough) Flush() {
}
