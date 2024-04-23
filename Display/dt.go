package Display

import (
	"errors"
	"sync"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	"github.com/gdamore/tcell"
)

// DtCell represents a cell in a table that can be drawn to the screen.
type DtCell struct {
	// Char represents the character of the cell.
	Char rune

	// Style represents the style of the cell.
	Style tcell.Style
}

func NewDtCell(char rune, style tcell.Style) DtCell {
	return DtCell{Char: char, Style: style}
}

// DrawTable represents a table that can be drawn to the screen.
type DrawTable struct {
	// width and height represent the width and height of the table, respectively.
	width, height int

	// cells represents the cells of the table.
	cells [][]DtCell

	// defaultStyle represents the default style of the table.
	defaultStyle tcell.Style

	// mu represents the mutex that will be used to lock the table.
	mu sync.RWMutex
}

// NewDrawTable creates a new table with the given width, height, and default style.
//
// If the width or height is less than 0, an error of type *ers.ErrInvalidParameter
// will be returned.
//
// Parameters:
//
//   - width: The width of the table.
//   - height: The height of the table.
//   - defaultStyle: The default style of the table.
//
// Returns:
//
//   - *DrawTable: A pointer to the new table.
//   - error: An error if the table could not be created.
func NewDrawTable(width, height int, defaultStyle tcell.Style) (*DrawTable, error) {
	if width < 0 {
		return nil, ers.NewErrInvalidParameter(
			"width",
			errors.New("value must be greater than or equal to 0"),
		)
	} else if height < 0 {
		return nil, ers.NewErrInvalidParameter(
			"height",
			errors.New("value must be greater than or equal to 0"),
		)
	}

	dt := &DrawTable{
		width:  width,
		height: height,
	}

	dt.cells = make([][]DtCell, height)
	for i := 0; i < height; i++ {
		dt.cells[i] = make([]DtCell, width)

		for j := 0; j < width; j++ {
			dt.cells[i][j] = NewDtCell(' ', defaultStyle)
		}
	}

	return dt, nil
}

// GetStyle returns the style of the table at the given position.
//
// If the position is out of bounds, an error of type *ers.ErrInvalidParameter
// will be returned.
//
// Parameters:
//
//   - x: The x-coordinate of the position.
//   - y: The y-coordinate of the position.
//   - value: The value to set.
//
// Returns:
//
//   - *DtCell: The style of the table at the given position.
//   - error: An error if the style could not be retrieved.
func (dt *DrawTable) WriteAt(x, y int, value *DtCell) error {
	if x < 0 || x >= dt.width {
		return ers.NewErrInvalidParameter(
			"x",
			ers.NewErrOutOfBounds(x, 0, dt.width),
		)
	} else if y < 0 || y >= dt.height {
		return ers.NewErrInvalidParameter(
			"y",
			ers.NewErrOutOfBounds(y, 0, dt.height),
		)
	}

	dt.mu.Lock()
	defer dt.mu.Unlock()

	if value == nil {
		dt.cells[y][x] = NewDtCell(' ', dt.defaultStyle)
	} else {
		dt.cells[y][x] = *value
	}

	return nil
}

// ResizeWidth resizes the width of the table. No op if the width is the same.
//
// If the width is less than 0, an error of type *ers.ErrInvalidParameter
// will be returned.
//
// Parameters:
//
//   - width: The new width of the table.
//
// Returns:
//
//   - error: An error if the width could not be resized.
func (dt *DrawTable) ResizeWidth(width int) error {
	if width == dt.width {
		return nil
	} else if width < 0 {
		return ers.NewErrInvalidParameter(
			"width",
			errors.New("value must be greater than or equal to 0"),
		)
	}

	dt.mu.Lock()
	defer dt.mu.Unlock()

	if width < dt.width {
		for i := 0; i < dt.height; i++ {
			dt.cells[i] = dt.cells[i][:width]
		}
	} else {
		for i := 0; i < dt.height; i++ {
			for j := dt.width; j < width; j++ {
				dt.cells[i] = append(dt.cells[i], NewDtCell(' ', dt.defaultStyle))
			}
		}
	}

	dt.width = width

	return nil
}

// ResizeHeight resizes the height of the table. No op if the height is the same.
//
// If the height is less than 0, an error of type *ers.ErrInvalidParameter
// will be returned.
//
// Parameters:
//
//   - height: The new height of the table.
//
// Returns:
//
//   - error: An error if the height could not be resized.
func (dt *DrawTable) ResizeHeight(height int) error {
	if height == dt.height {
		return nil
	} else if height < 0 {
		return ers.NewErrInvalidParameter(
			"height",
			errors.New("value must be greater than or equal to 0"),
		)
	}

	dt.mu.Lock()
	defer dt.mu.Unlock()

	if height < dt.height {
		dt.cells = dt.cells[:height]
	} else {
		for i := dt.height; i < height; i++ {
			row := make([]DtCell, dt.width)

			for j := 0; j < dt.width; j++ {
				row[j] = NewDtCell(' ', dt.defaultStyle)
			}

			dt.cells = append(dt.cells, row)
		}
	}

	dt.height = height

	return nil
}

// GetWidth returns the width of the table.
//
// Returns:
//
//   - int: The width of the table.
func (dt *DrawTable) GetWidth() int {
	return dt.width
}

// GetHeight returns the height of the table.
//
// Returns:
//
//   - int: The height of the table.
func (dt *DrawTable) GetHeight() int {
	return dt.height
}
