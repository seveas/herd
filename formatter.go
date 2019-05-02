package herd

import (
	"bytes"
	"fmt"
	"io"

	"github.com/mgutz/ansi"
)

var Formatters = map[string]Formatter{
	"pretty": PrettyFormatter{},
}

type Formatter interface {
	FormatHistoryItem(hi HistoryItem, w io.Writer)
	FormatCommand(c string, w io.Writer)
	FormatResult(r Result, w io.Writer)
}

type PrettyFormatter struct {
}

func (f PrettyFormatter) FormatHistoryItem(hi HistoryItem, w io.Writer) {
	f.FormatCommand(hi.Command, w)
	for _, h := range hi.Hosts {
		f.FormatResult(hi.Results[h.Name], w)
	}
}

func (f PrettyFormatter) FormatCommand(command string, w io.Writer) {
	fmt.Fprintln(w, ansi.Color(command, "cyan"))
}

func (f PrettyFormatter) FormatResult(r Result, w io.Writer) {
	if r.Err != nil {
		fmt.Fprintf(w, "%s%s %s%s\n", ansi.ColorCode("red"), r.Host, r.Err, ansi.ColorCode("reset"))
	} else {
		fmt.Fprintln(w, ansi.Color(r.Host, "green"))
	}
	if len(r.Stdout) > 0 {
		f.WriteIndented(w, r.Stdout)
	}
	if len(r.Stderr) != 0 {
		fmt.Fprintln(w, ansi.Color("----", "black+h"))
		f.WriteIndented(w, r.Stderr)
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
