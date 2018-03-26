package ast

import (
	"encoding/gob"
	"io"

	"github.com/frankreh/go-clang-v5.0/clang/cursorkind"
)

func (tu *TranslationUnit) EncodeGobV1(w io.Writer) error {
	enc := gob.NewEncoder(w) // Will write to w.
	/*
		-y- CursorKindId int // TBD
		-y- CursorNameId int // Id into CursorNameMap
		-n- Index        int // Own index into the TranslationUnit Cursors list.
		-y- ParentIndex  int //- 1 if Cursor is root.
		-n- Children     IndexPair
		-yy- Tokens       IndexPair
	*/
	// Deconstruct slice of Cursors into individual lists, omitting parts that
	// can be rebuild by a single traversal of the reconstructed slice later.
	ints := make([]int, len(tu.Cursors))

	// CursorKindId
	for i := range ints {
		ints[i] = int(tu.Cursors[i].CursorKindId)
	}
	if err := enc.Encode(ints); err != nil {
		return err
	}
	// CursorNameId
	for i := range ints {
		ints[i] = tu.Cursors[i].CursorNameId
	}
	if err := enc.Encode(ints); err != nil {
		return err
	}
	// ParentIndex
	for i := range ints {
		ints[i] = tu.Cursors[i].ParentIndex
	}
	if err := enc.Encode(ints); err != nil {
		return err
	}
	// TypeIndex
	for i := range ints {
		ints[i] = tu.Cursors[i].TypeIndex
	}
	if err := enc.Encode(ints); err != nil {
		return err
	}
	// Tokens.Head
	for i := range ints {
		ints[i] = tu.Cursors[i].Tokens.Head
	}
	if err := enc.Encode(ints); err != nil {
		return err
	}
	// Tokens.Len
	for i := range ints {
		ints[i] = tu.Cursors[i].Tokens.Len
	}
	if err := enc.Encode(ints); err != nil {
		return err
	}
	/*
		if err := enc.Encode(tu.Cursors); err != nil {
			return err
		}
	*/
	if err := enc.Encode(tu.TokenIds); err != nil {
		return err
	}
	if err := enc.Encode(tu.CursorNameMap); err != nil {
		return err
	}
	if err := enc.Encode(tu.TokenMap); err != nil {
		return err
	}
	if err := enc.Encode(tu.TokenNameMap); err != nil {
		return err
	}
	if err := enc.Encode(tu.TypeMap); err != nil {
		return err
	}
	if err := enc.Encode(tu.Back); err != nil {
		return err
	}
	return nil
}

func (tu *TranslationUnit) DecodeGobV1(r io.Reader) error {
	dec := gob.NewDecoder(r) // Will read from r.

	/*
		-y- CursorKindId int // TBD
		-y- CursorNameId int // Id into CursorNameMap
		-n- Index        int // Own index into the TranslationUnit Cursors list.
		-y- ParentIndex  int //- 1 if Cursor is root.
		-n- Children     IndexPair
		-yy- Tokens       IndexPair
	*/
	// Reconstruct slice of Cursors from individual lists, rebuidling parts that
	// can be rebuild by a single traversal of the reconstructed slice later.
	var ints []int

	// CursorKindId
	if err := dec.Decode(&ints); err != nil {
		return err
	}
	tu.Cursors = make([]Cursor, len(ints))
	for i := range ints {
		kind, err := cursorkind.Validate(ints[i])
		if err != nil {
			return err
		}
		tu.Cursors[i].CursorKindId = kind
	}
	ints = ints[:0]
	ints = nil
	// CursorNameId
	if err := dec.Decode(&ints); err != nil {
		return err
	}
	for i := range ints {
		tu.Cursors[i].CursorNameId = ints[i]
	}
	ints = ints[:0]
	ints = nil
	// ParentIndex
	if err := dec.Decode(&ints); err != nil {
		return err
	}
	for i := range ints {
		tu.Cursors[i].ParentIndex = ints[i]
		p := ints[i]
		if p < 0 {
			continue // Root will have parent set to -1.
		}
		c := &tu.Cursors[p].Children
		if c.Head == 0 {
			c.Head = i
		}
		c.Len += 1
	}
	ints = ints[:0]
	ints = nil
	// TypeIndex
	if err := dec.Decode(&ints); err != nil {
		return err
	}
	for i := range ints {
		tu.Cursors[i].TypeIndex = ints[i]
	}
	ints = ints[:0]
	ints = nil
	// Tokens.Head
	if err := dec.Decode(&ints); err != nil {
		return err
	}
	for i := range ints {
		tu.Cursors[i].Tokens.Head = ints[i]
	}
	ints = ints[:0]
	ints = nil
	// Tokens.Len
	if err := dec.Decode(&ints); err != nil {
		return err
	}
	for i := range ints {
		tu.Cursors[i].Tokens.Len = ints[i]
	}
	ints = ints[:0]
	ints = nil
	/*
		// All the code above to make encoding for Cursors more efficient.
		if err := dec.Decode(&tu.Cursors); err != nil {
			return err
		}
	*/
	if err := dec.Decode(&tu.TokenIds); err != nil {
		return err
	}
	if err := dec.Decode(&tu.CursorNameMap); err != nil {
		return err
	}
	if err := dec.Decode(&tu.TokenMap); err != nil {
		return err
	}
	if err := dec.Decode(&tu.TokenNameMap); err != nil {
		return err
	}
	if err := dec.Decode(&tu.TypeMap); err != nil {
		return err
	}
	if err := dec.Decode(&tu.Back); err != nil {
		return err
	}
	tu.DecodeFinish()
	return nil
}
