package Markdown

import (
	"errors"
	"strings"

	ers "github.com/PlayerR9/MyGoLibUnits/Errors"
)

const (
	// CodeBlockLimiter is the string used to denote the start and end of a code block.
	CodeBlockLimiter string = "```"
)

// CommonLanguageName represents the common language names that can be used in a code block.
type CommonLanguageName int

const (
	// PlainText represents the plaintext language.
	PlainText CommonLanguageName = iota

	// OtherCML represents any non-standard language.
	OtherCML
)

// String is a method of fmt.Stringer interface that returns the string representation
// of the language name.
//
// Returns:
//
//   - string: The string representation of the language name.
func (cln CommonLanguageName) String() string {
	return [...]string{
		"plaintext",
	}[cln]
}

// ParseCommonLanguageName parses a common language name from a string.
//
// Parameters:
//
//   - name: The name of the language.
//
// Returns:
//
//   - CommonLanguageName: The common language name. If the name is not recognized,
//     -1 is returned.
func ParseCommonLanguageName(name string) CommonLanguageName {
	switch strings.ToLower(name) {
	case "plaintext":
		return PlainText
	default:
		return OtherCML
	}
}

// CodeBlock represents a code block in a markdown file.
type CodeBlock struct {
	// Language represents the language of the code block.
	Language string

	// Content represents the content of the code block.
	Content []string
}

// ExtractCodeBlockAt extracts a code block at the given position and returns
// the code block and the new position.
//
// Parameters:
//
//   - lines: A slice of strings containing the lines of the file.
//   - at: The position to start parsing.
//
// Returns:
//
//   - CodeBlock: The code block.
//   - int: The new position (this is exclusive, i.e. the next line to parse).
//   - error: An error if the code block is not formatted correctly.
func ExtractCodeBlockAt(lines []string, at int) (CodeBlock, int, error) {
	cb := CodeBlock{
		Content: make([]string, 0),
	}

	if at >= len(lines) {
		return cb, at, ers.NewErrInvalidParameter(
			"at",
			ers.NewErrOutOfBounds(at, 0, len(lines)),
		)
	} else if len(lines) < 2 {
		return cb, at, ers.NewErrInvalidParameter(
			"lines",
			errors.New("not enough lines to extract code block"),
		)
	}

	if !strings.HasPrefix(lines[at], CodeBlockLimiter) {
		return cb, at, errors.New("line does not start with code block limiter")
	}

	cb.Language = strings.TrimPrefix(lines[at], CodeBlockLimiter)

	for i := at + 1; i < len(lines); i++ {
		if lines[i] == CodeBlockLimiter {
			return cb, i + 1, nil
		}

		cb.Content = append(cb.Content, lines[i])
	}

	return cb, at, errors.New("no code block end found")
}
