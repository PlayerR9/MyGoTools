package Lexer

import (
	gr "github.com/PlayerR9/MyGoLib/Utility/Grammar"
)

// TokenStatus represents the status of a token
type TokenStatus int

const (
	// TkComplete represents a token that has been fully lexed
	TkComplete TokenStatus = iota

	// TkIncomplete represents a token that has not been fully lexed
	TkIncomplete

	// TkError represents a token that has an error
	TkError
)

// String returns the string representation of a TokenStatus
//
// Returns:
//
//   - string: The string representation of the TokenStatus
func (s TokenStatus) String() string {
	return [...]string{
		"complete",
		"incomplete",
		"error",
	}[s]
}

// helperToken is a wrapper around a *LeafToken that adds a status field
type helperToken struct {
	// Status is the status of the token
	Status TokenStatus

	// Tok is the *LeafToken
	Tok *gr.LeafToken
}

// SetStatus sets the status of the token
//
// Parameters:
//
//   - status: The status to set
func (ht *helperToken) SetStatus(status TokenStatus) {
	ht.Status = status
}
