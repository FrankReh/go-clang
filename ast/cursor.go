package ast

import "github.com/frankreh/go-clang/clang/cursorkind"

type IndexPair struct {
	Head int
	Len  int
}

// Next returns the next index after the list that this pair of values represents.
func (ip IndexPair) Next() int {
	return ip.Head + ip.Len
}

// Cursor is the pure Go version of the clang Cursor.
// It exists as part of the TranslationUniit Cursors list.
type Cursor struct {
	CursorKindId cursorkind.Kind
	CursorNameId int // Id into CursorNameMap
	ParentIndex  int //- 1 if Cursor is root.
	TypeIndex    int // Index into TypeSlice
	Children     IndexPair
	Tokens       IndexPair

	// Children does not have to be initialized within Visit and does not have
	// to be serialized. It can be recomputed from the Cursor position in the
	// overall list, and the ParentIndex.

}
