package Lexer

import (
	"errors"
	"slices"

	gr "github.com/PlayerR9/MyGoLib/Utility/Grammar"

	nd "github.com/PlayerR9/MyGoLib/CustomData/Node"
	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

type TokenStatus int

const (
	TkComplete TokenStatus = iota
	TkIncomplete
	TkError
)

func (s TokenStatus) String() string {
	return [...]string{
		"complete",
		"incomplete",
		"error",
	}[s]
}

type helperToken struct {
	Status TokenStatus
	Tok    *gr.LeafToken
}

func (ht *helperToken) SetStatus(status TokenStatus) {
	ht.Status = status
}

func (ht *helperToken) GetStatus() TokenStatus {
	return ht.Status
}

type Lexer struct {
	grammar *gr.Grammar
	toSkip  []string
	root    *nd.Node[*helperToken]
	leaves  []*nd.Node[*helperToken]
}

func NewLexer(grammar *gr.Grammar) (*Lexer, error) {
	if grammar == nil {
		return nil, ers.NewErrInvalidParameter("grammar").
			Wrap(errors.New("grammar cannot be nil"))
	}

	lex := &Lexer{
		grammar: grammar,
		toSkip:  grammar.LhsToSkip,
	}

	return lex, nil
}

func (l *Lexer) addFirstLeaves(b []byte) error {
	matches := l.grammar.Match(0, b)
	if len(matches) == 0 {
		return errors.New("no matches found at index 0")
	}

	// Get the longest match
	matches = getLongestMatches(matches)
	for _, match := range matches {
		leafToken, ok := match.Matched.(*gr.LeafToken)
		if !ok {
			return errors.New("this should not happen: match.Matched is not a *LeafToken")
		}

		l.root.AddChild(&helperToken{
			Status: TkIncomplete,
			Tok:    leafToken,
		})
		l.leaves = l.root.GetLeaves()
	}

	return nil
}

func (l *Lexer) processLeaf(leaf *nd.Node[*helperToken], b []byte) {
	nextAt := leaf.Data.Tok.GetPos() + len(leaf.Data.Tok.GetData().(string))
	if nextAt >= len(b) {
		leaf.Data.SetStatus(TkComplete)
		return
	}
	subset := b[nextAt:]

	matches := l.grammar.Match(nextAt, subset)

	if len(matches) == 0 {
		// Branch is done but no match found
		leaf.Data.SetStatus(TkError)
		return
	}

	// Get the longest match
	matches = getLongestMatches(matches)
	for _, match := range matches {
		leafToken, ok := match.Matched.(*gr.LeafToken)
		if !ok {
			leaf.Data.SetStatus(TkError)
			return
		}

		leaf.AddChild(&helperToken{
			Status: TkIncomplete,
			Tok:    leafToken,
		})
	}

	leaf.Data.SetStatus(TkComplete)
}

func getLongestMatches(matches []gr.MatchedResult) []gr.MatchedResult {
	return slext.FilterByPositiveWeight(matches, func(match gr.MatchedResult) (int, bool) {
		val, ok := match.Matched.GetData().(string)
		if !ok {
			return 0, false
		}

		return len(val), true
	})
}

func (l *Lexer) Lex(b []byte) error {
	if len(b) == 0 {
		return errors.New("no tokens to parse")
	}

	l.root = nd.NewNode(&helperToken{
		Status: TkIncomplete,
		Tok:    gr.NewLeafToken("root", "", -1),
	})

	if err := l.addFirstLeaves(b); err != nil {
		return err
	}

	l.root.Data.SetStatus(TkComplete)

	for {
		// Remove all the leaves that are completed
		todo := slext.SliceFilter(l.leaves, func(leaf *nd.Node[*helperToken]) bool {
			return leaf.Data.GetStatus() != TkComplete
		})
		if len(todo) == 0 {
			// All leaves are complete
			break
		}

		// Remove all the leaves that are in error
		todo = slext.SliceFilter(todo, func(leaf *nd.Node[*helperToken]) bool {
			return leaf.Data.GetStatus() != TkError
		})
		if len(todo) == 0 {
			// All leaves are in error
			break
		}

		// Remaining leaves are incomplete
		var newLeaves []*nd.Node[*helperToken]

		for _, leaf := range todo {
			l.processLeaf(leaf, b)
			newLeaves = append(newLeaves, leaf.GetLeaves()...)
		}

		l.leaves = newLeaves
	}

	return nil
}

func (l *Lexer) removeSkippedTokens(tokenBranches [][]*gr.LeafToken) [][]*gr.LeafToken {
	// Remove the root token
	for i, branch := range tokenBranches {
		tokenBranches[i] = branch[1:]
	}

	for i, branch := range tokenBranches {
		tokenBranches[i] = slext.SliceFilter(branch, func(token *gr.LeafToken) bool {
			return !slices.Contains(l.toSkip, token.GetID())
		})
	}

	return tokenBranches
}

func completeBranchFilter(tokens []*helperToken) bool {
	return !slices.ContainsFunc(tokens, func(token *helperToken) bool {
		return token.GetStatus() != TkComplete
	})
}

func emptyBranchFilter(tokens []*gr.LeafToken) bool {
	return len(tokens) > 0
}

func (l *Lexer) GetTokens() ([][]*gr.LeafToken, error) {
	if l.root == nil {
		return nil, errors.New("must call Lexer.Lex() first")
	}

	tokenBranches := l.root.SnakeTraversal()
	tokenBranches = slext.SliceFilter(tokenBranches, completeBranchFilter)

	// Convert the tokens to []*LeafToken
	result := make([][]*gr.LeafToken, len(tokenBranches))

	for i, branch := range tokenBranches {
		result[i] = make([]*gr.LeafToken, len(branch))

		for j, token := range branch {
			result[i][j] = token.Tok
		}
	}

	result = l.removeSkippedTokens(result)
	result = slext.SliceFilter(result, emptyBranchFilter)

	return result, nil
}

func FullLexer(grammar *gr.Grammar, content string) ([][]*gr.LeafToken, error) {
	lexer, err := NewLexer(grammar)
	if err != nil {
		return nil, err
	}

	err = lexer.Lex([]byte(content))
	tokens, _ := lexer.GetTokens()

	return tokens, err
}
