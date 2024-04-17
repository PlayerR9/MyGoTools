package Display

import "errors"

type ErrESCPressed struct{}

func (e ErrESCPressed) Error() string {
	return "ESC key pressed"
}

func NewErrESCPressed() *ErrESCPressed {
	return &ErrESCPressed{}
}

func IsESCPressed(err error) bool {
	var errESC *ErrESCPressed

	return errors.As(err, &errESC)
}
