package Display

type WriteOnlyDTer interface {
	WriteAt(int, int, *DtCell) error
	GetWidth() int
	GetHeight() int
}
