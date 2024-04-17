package Display

import "github.com/gdamore/tcell"

type DtCell struct {
	Char  rune
	Style tcell.Style
}

func NewDtCell(char rune, style tcell.Style) DtCell {
	return DtCell{Char: char, Style: style}
}
