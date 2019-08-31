package herd

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/mgutz/ansi"
)

var Formatters = map[string]Formatter{
	"pretty": PrettyFormatter{},
}

type Formatter interface {
	FormatCommand(c string, w io.Writer)
	FormatResult(r Result, w io.Writer)
	FormatStatus(r Result, w io.Writer)
}

type PrettyFormatter struct {
}

func (f PrettyFormatter) FormatCommand(command string, w io.Writer) {
	fmt.Fprintln(w, ansi.Color(command, "cyan"))
}

func (f PrettyFormatter) FormatResult(r Result, w io.Writer) {
	if r.Err != nil {
		fmt.Fprintf(w, ansi.Color(r.Host, "red")+" ")
		f.FormatStatus(r, w)
	} else {
		fmt.Fprintf(w, ansi.Color(r.Host, "green")+" ")
		f.FormatStatus(r, w)
	}
	if len(r.Stdout) > 0 {
		f.WriteIndented(w, r.Stdout)
	}
	if len(r.Stderr) != 0 {
		fmt.Fprintln(w, ansi.Color("----", "black+h"))
		f.WriteIndented(w, r.Stderr)
	}
}

func (f PrettyFormatter) FormatStatus(r Result, w io.Writer) {
	if r.Err != nil {
		fmt.Fprintln(w, ansi.Color(fmt.Sprintf("%s after %s", r.Err, r.EndTime.Sub(r.StartTime).Truncate(time.Second)), "red"))
	} else {
		fmt.Fprintln(w, ansi.Color(fmt.Sprintf("completed successfully after %s", r.EndTime.Sub(r.StartTime).Truncate(time.Second)), "green"))
	}
}

func (f PrettyFormatter) WriteIndented(w io.Writer, msg []byte) {
	w.Write([]byte{0x20, 0x20, 0x20, 0x20})
	if msg[len(msg)-1] == 0x0a {
		msg = msg[:len(msg)-1]
	}
	w.Write(bytes.ReplaceAll(msg, []byte{0x0a}, []byte{0x0a, 0x20, 0x20, 0x20, 0x20}))
	w.Write([]byte{0x0a})
}
