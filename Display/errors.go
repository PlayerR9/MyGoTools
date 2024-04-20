package Display

import "errors"

// ErrESCPressed is an error that is returned when the ESC key is pressed.
type ErrESCPressed struct{}

// Error returns the error message.
//
// Returns:
//
//   - string: The error message.
func (e ErrESCPressed) Error() string {
	return "ESC key pressed"
}

// NewErrESCPressed creates a new ErrESCPressed error.
//
// Returns:
//
//   - *ErrESCPressed: A pointer to the new error.
func NewErrESCPressed() *ErrESCPressed {
	return &ErrESCPressed{}
}

// IsESCPressed checks if the given error is of type ErrESCPressed.
//
// Parameters:
//
//   - err: The error to check.
//
// Returns:
//
//   - bool: True if the error is of type ErrESCPressed, false otherwise.
func IsESCPressed(err error) bool {
	var errESC *ErrESCPressed

	return errors.As(err, &errESC)
}
