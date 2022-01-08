package herd

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type datawriter interface {
	Write([]string) error
	Flush()
}

type columnizer struct {
	width   int
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
	tl := len(c.lengths) - 1
	for i, v := range r {
		if l := len(v); l > c.lengths[i] {
			c.lengths[i] = l
		}
		tl += c.lengths[i]
	}
	c.width = tl
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

type ByteWriter interface {
	Write([]byte) (int, error)
	Bytes() []byte
}

type lineWriterBuffer struct {
	oc      chan OutputLine
	host    *Host
	stderr  bool
	buf     *bytes.Buffer
	lineBuf []byte
}

func NewLineWriterBuffer(host *Host, stderr bool, oc chan OutputLine) *lineWriterBuffer {
	return &lineWriterBuffer{
		buf:     bytes.NewBuffer([]byte{}),
		lineBuf: []byte{},
		host:    host,
		oc:      oc,
		stderr:  stderr,
	}
}

func (buf *lineWriterBuffer) Write(p []byte) (int, error) {
	n, err := buf.buf.Write(p)
	buf.lineBuf = bytes.Join([][]byte{buf.lineBuf, p}, []byte{})
	for {
		idx := bytes.Index(buf.lineBuf, []byte("\n"))
		if idx == -1 {
			break
		}
		buf.oc <- OutputLine{Host: buf.host, Data: buf.lineBuf[:idx+1], Stderr: buf.stderr}
		buf.lineBuf = buf.lineBuf[idx+1:]
	}
	return n, err
}

func (buf *lineWriterBuffer) Bytes() []byte {
	return buf.buf.Bytes()
}
