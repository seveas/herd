package katyusha

import (
	"io"
)

type Formatter interface {
	Format(hi HistoryItem, w io.Writer)
	FormatResult(r Result, w io.Writer)
}
