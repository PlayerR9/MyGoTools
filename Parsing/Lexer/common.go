package Lexer

import gr "github.com/PlayerR9/MyGoLib/Utility/Grammar"

// LexString is a function that, given an input string, returns a slice of tokens.
//
// Parameters:
//
//   - input: The input string.
//
// Returns:
//
//   - [][]gr.LeafToken: A slice of slices of tokens.
//   - error: An error if the input string cannot be lexed.
func LexString(lexer *Lexer, input string) ([][]gr.LeafToken, error) {
	err := lexer.Lex([]byte(input))
	if err != nil {
		return nil, err
	}

	return lexer.GetTokens()
}
