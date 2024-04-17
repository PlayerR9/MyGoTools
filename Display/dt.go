package Display

import (
	"errors"
	"sync"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	"github.com/gdamore/tcell"
)

type DrawTable struct {
	width, height int
	cells         [][]DtCell

	defaultStyle tcell.Style

	mu sync.RWMutex
}

func NewDrawTable(width, height int, defaultStyle tcell.Style) *DrawTable {
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

	return dt
}

func (dt *DrawTable) WriteAt(x, y int, value *DtCell) error {
	if x < 0 || x >= dt.width {
		return ers.NewErrInvalidParameter("x").
			Wrap(ers.NewErrOutOfBound(x, 0, dt.width))
	} else if y < 0 || y >= dt.height {
		return ers.NewErrInvalidParameter("y").
			Wrap(ers.NewErrOutOfBound(y, 0, dt.height))
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

func (dt *DrawTable) ResizeWidth(width int) error {
	if width == dt.width {
		return nil
	} else if width < 0 {
		return ers.NewErrInvalidParameter("width").
			Wrap(errors.New("value must be greater than or equal to 0"))
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

func (dt *DrawTable) ResizeHeight(height int) error {
	if height == dt.height {
		return nil
	} else if height < 0 {
		return ers.NewErrInvalidParameter("height").
			Wrap(errors.New("value must be greater than or equal to 0"))
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

func (dt *DrawTable) GetWidth() int {
	return dt.width
}

func (dt *DrawTable) GetHeight() int {
	return dt.height
}
