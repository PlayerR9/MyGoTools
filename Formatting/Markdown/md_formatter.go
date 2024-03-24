// Unfinished implementation of a markdown formatter
package Markdown

import (
	"fmt"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

type TableAllignmentType int

const (
	TableStyleCenter TableAllignmentType = iota
	TableStyleLeft
	TableStyleRight
)

func (a TableAllignmentType) String() string {
	return [...]string{
		":---:",
		":---",
		"---:",
	}[a]
}

var (
	tokPipe      string = " | "
	tokNewLine   string = "\n"
	PedanticMode bool   = true
)

type Header struct {
	Text  string
	Style TableAllignmentType
}

type Table struct {
	headers []Header
	rows    [][]string
}

func NewTable() *Table {
	return &Table{
		headers: make([]Header, 0),
		rows:    make([][]string, 0),
	}
}

func (t *Table) WithHeader(header Header) *Table {
	t.AppendHeader(header)

	return t
}

func (t *Table) AppendHeader(header Header) error {
	t.headers = append(t.headers, header)

	if !PedanticMode {
		for i, row := range t.rows {
			if len(row) < len(t.headers) {
				return fmt.Errorf("row %d has less elements than the number of headers", i)
			}
		}
	} else {
		for _, row := range t.rows {
			if len(row) >= len(t.headers) {
				continue
			}

			for len(row) < len(t.headers) {
				row = append(row, "")
			}
		}
	}

	return nil
}

func (t *Table) CanAddRow(elements []string) bool {
	return len(elements) == len(t.headers)
}

func (t *Table) AddRow(elements []string) error {
	if len(elements) == len(t.headers) {
		t.rows = append(t.rows, elements)
		return nil
	}

	return ers.NewErrInvalidParameter("elements").
		Wrap(fmt.Errorf("number of elements (%d) does not match number of headers (%d)",
			len(elements), len(t.headers)))
}

func rowToString[T any](row []T, f func(T) string) string {
	var builder strings.Builder

	builder.WriteString(tokPipe)

	for _, element := range row {
		fmt.Fprintf(&builder, "%s%s", f(element), tokPipe)
	}

	builder.WriteString(tokNewLine)

	return builder.String()
}

func (t *Table) ToText() []string {
	text := make([]string, 0, 2*len(t.headers)+len(t.rows)+1)

	text = append(text, rowToString[Header](t.headers, func(h Header) string {
		return h.Text
	}))
	text = append(text, rowToString[Header](t.headers, func(h Header) string {
		return h.Style.String()
	}))

	for _, row := range t.rows {
		text = append(text, rowToString(row, func(s string) string {
			return s
		}))
	}

	return text
}
