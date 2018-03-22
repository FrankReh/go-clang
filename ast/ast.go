// Package ast implements major parts of the abstract syntax tree that the libclang creates.
//
// The ast package does not use cgo nor does it use the unsafe package.
//
// The primary data structures in this package are also pointer free, and are expressly designed
// for efficient serialization.
package ast

import (
	"fmt"
	"strings"
)

// --- 1803 clang ast
// Idea is to create one list of all tokens in tu
// and one list of all cursors in tu. Actually two lists, that mirror each other,
// one of actual clang.cursors and the other of the other of main package Cursors
// that will contain all the computed stuff that can be stored and reloaded without cgo.

// TranslationUnit is the pure Go version of the clang translation unit.
// Not using pointers so it can be serialized.
type TranslationUnit struct {
	Cursors  []Cursor  // [0] is the Root Cursor.
	TokenIds []TokenId // TokenId is mapped by TokenMap to get a Token.

	// String/Id maps that allow the Cursor and Token to use int IDs properties.
	// TBD remove CursorKindMap StringMap
	CursorNameMap StringMap
	TokenMap      TokenMap // Token:"{TokenKindId:2 TokenNameId:3}" mapped to ID
	// TBD remove TokenKindMap  StringMap
	TokenNameMap StringMap
}

// DecodeFinish Completes the setup of the lists and maps after the object
// has been decoded.
func (tu *TranslationUnit) DecodeFinish() {
	//tu.CursorKindMap.DecodeFinish()
	tu.CursorNameMap.DecodeFinish()
	tu.TokenMap.DecodeFinish()
	//tu.TokenKindMap.DecodeFinish()
	tu.TokenNameMap.DecodeFinish()
}

func (tu *TranslationUnit) AssertEqual(tu2 *TranslationUnit) error {
	if tu == tu2 {
		return nil
	}
	/*
		if err := tu.CursorKindMap.AssertEqual(&tu2.CursorKindMap); err != nil {
			return err
		}
	*/
	if err := tu.CursorNameMap.AssertEqual(&tu2.CursorNameMap); err != nil {
		return err
	}
	if err := tu.TokenMap.AssertEqual(&tu2.TokenMap); err != nil {
		return err
	}
	/*
		if err := tu.TokenKindMap.AssertEqual(&tu2.TokenKindMap); err != nil {
			return err
		}
	*/
	if err := tu.TokenNameMap.AssertEqual(&tu2.TokenNameMap); err != nil {
		return err
	}
	return nil
}

// numberStrings return string of lines with each slice element shown with its
// slice position prefixed. Before letting line get past width, insert a newline.
func numberStrings(list []string, width int) string {
	b := new(strings.Builder)
	line := 0
	for i, s := range list {
		str := fmt.Sprintf("%d:%s", i, s)
		if 1+line+len(str) > width {
			fmt.Fprintf(b, "\n")
			line = 0
		}
		if line > 0 {
			fmt.Fprintf(b, " ")
			line += 1
		}
		fmt.Fprintf(b, "%s", str)
		line += len(str)
	}

	return b.String()
}

func (t TranslationUnit) GoString() string {
	// Convert two slices to slices of strings
	tokens := make([]string, len(t.TokenIds))
	for i := range tokens {
		tokens[i] = fmt.Sprintf("%v", t.TokenIds[i])
	}
	tm_tokens := make([]string, len(t.TokenMap.Tokens))
	for i := range tm_tokens {
		tm_tokens[i] = fmt.Sprintf("%v", t.TokenMap.Tokens[i])
	}
	cursors := make([]string, len(t.Cursors))
	for i := range cursors {
		cursors[i] = fmt.Sprintf("%v", t.Cursors[i])
	}
	width := 120
	b := new(strings.Builder)
	fmt.Fprintf(b, "Tokens:\n%v\n", numberStrings(tokens, width))
	fmt.Fprintf(b, "TokenMap:\n%v\n", numberStrings(tm_tokens, width))
	//fmt.Fprintf(b, "TokenKindMap:\n%v\n", numberStrings(t.TokenKindMap.Strings, width))
	fmt.Fprintf(b, "TokenNameMap:\n%v\n", numberStrings(t.TokenNameMap.Strings, width))
	fmt.Fprintf(b, "Cursors:\n%v\n", numberStrings(cursors, width))
	//fmt.Fprintf(b, "CursorKindMap:\n%v\n", numberStrings(t.CursorKindMap.Strings, width))
	fmt.Fprintf(b, "CursorNameMap:\n%v\n", numberStrings(t.CursorNameMap.Strings, width))
	return b.String()
}
