package Display

// WriteOnlyDTer is an interface that represents a table that can be written to.
type WriteOnlyDTer interface {
	// WriteAt writes a cell to the table at the given x and y coordinates.
	//
	// If the x-coordinate is out of bounds, an error of type *ers.ErrInvalidParameter
	// will be returned.
	//
	// If the y-coordinate is out of bounds, an error of type *ers.ErrInvalidParameter
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
	//   - error: An error if the cell could not be written.
	WriteAt(x, y int, value *DtCell) error

	// GetWidth returns the width of the table.
	//
	// Returns:
	//
	//   - int: The width of the table.
	GetWidth() int

	// GetHeight returns the height of the table.
	//
	// Returns:
	//
	//   - int: The height of the table.
	GetHeight() int
}
