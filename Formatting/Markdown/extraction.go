package Markdown

import (
	"errors"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

const (
	// SectionStyling is the styling used to denote a section in the markdown file.
	SectionStyling string = "***"
)

// ExtractSectionAt extracts a section at the given position and returns
// the section content and the new position.
//
// Parameters:
//
//   - lines: A slice of strings containing the lines of the file.
//   - at: The position to start parsing.
//
// Returns:
//
//   - string: The section content.
//   - int: The new position (this is exclusive, i.e. the next line to parse).
//   - error: An error if the section is not formatted correctly.
func ExtractSectionAt(lines []string, at int) (string, int, error) {
	if at >= len(lines) {
		return "", at, ers.NewErrInvalidParameter(
			"at",
			ers.NewErrOutOfBounds(at, 0, len(lines)),
		)
	}

	if !strings.HasPrefix(lines[at], SectionStyling) {
		return "", at, errors.New("missing section styling at the beginning")
	} else if !strings.HasSuffix(lines[at], SectionStyling) {
		return "", at, errors.New("missing section styling at the end")
	}

	section := strings.TrimPrefix(lines[at], SectionStyling)
	section = strings.TrimSuffix(section, SectionStyling)

	return section, at + 1, nil
}
