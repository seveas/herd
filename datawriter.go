package herd

import (
	"fmt"
	"io"
)

type datawriter interface {
	Write([]string) error
	Flush()
}

type Columnizer struct {
	rows    [][]string
	lengths []int
	output  io.Writer
	sep     string
}

func NewColumnizer(w io.Writer, sep string) *Columnizer {
	return &Columnizer{rows: make([][]string, 0), output: w, sep: sep}
}

func (c *Columnizer) Write(r []string) error {
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

func (c *Columnizer) Flush() {
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
