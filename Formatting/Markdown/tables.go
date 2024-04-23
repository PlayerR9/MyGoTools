// Unfinished implementation of a markdown formatter
package Markdown

import (
	"fmt"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// TableAllignmentType represents the type of alignment that a table can have.
type TableAllignmentType int8

const (
	// TableStyleCenter represents a center alignment.
	TableStyleCenter TableAllignmentType = iota

	// TableStyleLeft represents a left alignment.
	TableStyleLeft

	// TableStyleRight represents a right alignment.
	TableStyleRight
)

// String is a methos of fmt.Stringer interface that returns the string representation of the alignment.
//
// Returns:
//
//   - string: The string representation of the alignment.
func (a TableAllignmentType) String() string {
	return [...]string{
		":---:",
		":---",
		"---:",
	}[a]
}

const (
	// TokPipe is the pipe character.
	TokPipe string = " | "

	// TokNewLine is the newline character.
	TokNewLine string = "\n"
)

var (
	// PedanticMode is a flag that determines if the formatter should be pedantic.
	//
	// If the formatter is pedantic, it will add empty elements to rows that do not have
	// the same number of elements as the number of headers.
	//
	// If the formatter is not pedantic, it will return an error if a row does not have
	// the same number of elements as the number of headers.
	PedanticMode bool = true
)

// Header represents a header in a table.
type Header struct {
	// Text represents the text of the header.
	Text string

	// Style represents the style of the header.
	Style TableAllignmentType
}

// Table represents a table.
type Table struct {
	// headers represents the headers of the table.
	headers []Header

	// rows represents the rows of the table.
	rows [][]string
}

// NewTable creates a new table.
//
// Returns:
//
//   - *Table: A pointer to the new table.
func NewTable() *Table {
	return &Table{
		headers: make([]Header, 0),
		rows:    make([][]string, 0),
	}
}

// WithHeader adds a header to the table.
//
// Parameters:
//
//   - header: The header to add.
//
// Returns:
//
//   - *Table: A pointer to the table.
func (t *Table) WithHeader(header Header) *Table {
	t.AppendHeader(header)

	return t
}

// AppendHeader appends a header to the table.
//
// Parameters:
//
//   - header: The header to append.
//
// Returns:
//
//   - error: An error if the header could not be appended.
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

// CanAddRow returns true if the number of elements in the row matches the number of headers.
//
// Parameters:
//
//   - elements: The elements of the row.
//
// Returns:
//
//   - bool: True if the number of elements in the row matches the number of headers.
func (t *Table) CanAddRow(elements []string) bool {
	return len(elements) == len(t.headers)
}

// AddRow adds a row to the table.
//
// Parameters:
//
//   - elements: The elements of the row.
//
// Returns:
//
//   - error: An error if the row could not be added.
func (t *Table) AddRow(elements []string) error {
	if len(elements) == len(t.headers) {
		t.rows = append(t.rows, elements)
		return nil
	}

	return ers.NewErrInvalidParameter(
		"elements",
		fmt.Errorf("number of elements (%d) does not match number of headers (%d)",
			len(elements), len(t.headers)),
	)
}

// rowToString converts a row to a string.
//
// Parameters:
//
//   - row: The row to convert.
//   - f: The function that converts an element to a string.
//
// Returns:
//
//   - string: The string representation of the row.
func rowToString[T any](row []T, f func(T) string) string {
	var builder strings.Builder

	builder.WriteString(TokPipe)

	for _, element := range row {
		fmt.Fprintf(&builder, "%s%s", f(element), TokPipe)
	}

	builder.WriteString(TokNewLine)

	return builder.String()
}

// ToText converts the table to slice of strings.
//
// Returns:
//
//   - []string: The string representation of the table.
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
