package Lexer

import (
	gr "github.com/PlayerR9/MyGoLib/Utility/Grammar"

	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// getLongestMatches returns the longest matches
//
// Parameters:
//
//   - matches: The matches to filter
//
// Returns:
//
//   - []MatchedResult: The longest matches
func getLongestMatches(matches []gr.MatchedResult) []gr.MatchedResult {
	return slext.FilterByPositiveWeight(matches, func(match gr.MatchedResult) (int, bool) {
		leaf, ok := match.Matched.(*gr.LeafToken)
		if !ok {
			return 0, false
		}

		return len(leaf.Data), true
	})
}

// emptyBranchFilter is a filter function that returns true if a branch is not empty
//
// Parameters:
//
//   - tokens: The tokens to check
//
// Returns:
//
//   - bool: True if the branch is not empty, false otherwise
func emptyBranchFilter(tokens []gr.LeafToken) bool {
	return len(tokens) > 0
}

func filterInvalidBranches(branches [][]helperToken) ([][]helperToken, int) {
	branches, ok := slext.SFSeparateEarly(branches, func(h []helperToken) bool {
		return len(h) != 0 && h[len(h)-1].Status == TkComplete
	})
	if ok {
		return branches, -1
	}

	// Return the longest branch
	branches = slext.FilterByPositiveWeight(branches, func(h []helperToken) (int, bool) {
		return len(h), true
	})

	return [][]helperToken{branches[0]}, len(branches[0])
}
