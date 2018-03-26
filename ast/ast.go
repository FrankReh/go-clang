// Package ast implements major parts of the abstract syntax tree that the libclang creates.
//
// The ast package does not use cgo nor does it use the unsafe package.
//
// The primary data structures in this package are also pointer free, and are expressly designed
// for efficient serialization.
package ast

import (
	"fmt"
	"sort"
	"strings"

	"github.com/frankreh/go-clang-v5.0/clang/cursorkind"
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
	CursorNameMap StringMap
	TokenMap      TokenMap // Token:"{TokenKindId:2 TokenNameId:3}" mapped to ID
	TokenNameMap  StringMap
	TypeMap       TypeMap // Type:{TypeKind Index} Index to specific slice, given the TypeKind.

	// Create back reference "pointers" to recognize with cursor tree points back
	// to itself.
	Back map[int]int
}

// DecodeFinish Completes the setup of the lists and maps after the object
// has been decoded.
func (tu *TranslationUnit) DecodeFinish() {
	tu.CursorNameMap.DecodeFinish()
	tu.TokenMap.DecodeFinish()
	tu.TokenNameMap.DecodeFinish()
	//Maybe don't call this here after all.
	//tu.SetBackChildren()
}

// SetBackChildren sets the Children (probably one child) for
// each of the Back Cursors. Done by looking through the Back map,
// since looking through the Cursors for a Back kind would not lead
// to knowing where to point anywhere.
func (tu *TranslationUnit) SetBackChildren() {
	for backId, seenId := range tu.Back {
		if backId <= seenId {
			// TBD make these errors
			panic("back not after seen")
		}
		if backId >= len(tu.Cursors) {
			panic("back out of range")
		}
		if seenId >= len(tu.Cursors) {
			panic("seen out of range")
		}
		if tu.Cursors[backId].CursorKindId != cursorkind.Back {
			panic("backId not leading to Back cursor")
		}
		if tu.Cursors[backId].Children.Len != 0 {
			panic("backId leads to Back cursor with children already")
		}
		tu.Cursors[backId].Children.Head = seenId
		tu.Cursors[backId].Children.Len++
	}
}

func (tu *TranslationUnit) AssertEqual(tu2 *TranslationUnit) error {
	if tu == tu2 {
		return nil
	}
	if err := tu.CursorNameMap.AssertEqual(&tu2.CursorNameMap); err != nil {
		return err
	}
	if err := tu.TokenMap.AssertEqual(&tu2.TokenMap); err != nil {
		return err
	}
	if err := tu.TokenNameMap.AssertEqual(&tu2.TokenNameMap); err != nil {
		return err
	}
	if err := tu.TypeMap.AssertEqual(&tu2.TypeMap); err != nil {
		return err
	}
	if err := assertEqualMaps(tu.Back, tu2.Back); err != nil {
		return err
	}
	return nil
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
	fmt.Fprintf(b, "TokenNameMap:\n%v\n", numberStrings(t.TokenNameMap.Strings, width))
	fmt.Fprintf(b, "TypeMap:\n%#v\n", t.TypeMap)
	fmt.Fprintf(b, "Cursors:\n%v\n", numberStrings(cursors, width))
	fmt.Fprintf(b, "CursorNameMap:\n%v\n", numberStrings(t.CursorNameMap.Strings, width))
	if len(t.Back) > 0 {
		fmt.Fprintf(b, "Back:\n%v\n", numberStringsNoIndex(backStrings(t.Back), width))
	}
	return b.String()
}

// backStrings returns strings of form "key:value", sorted by key.
// Sort them because their display gets compared to an expected value.
func backStrings(m map[int]int) []string {
	// To store the keys in slice in sorted order
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	r := make([]string, len(keys))
	for i, k := range keys {
		r[i] = fmt.Sprintf("%d:%d", k, m[k])
	}
	return r
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

// numberStringsNoIndex return string of lines with each slice element shown.
// Before letting line get past width, insert a newline.
func numberStringsNoIndex(list []string, width int) string {
	b := new(strings.Builder)
	line := 0
	for _, s := range list {
		str := s
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
func assertEqualMaps(a, b map[int]int) error {
	if len(a) != len(b) {
		return fmt.Errorf("map lengths a %d, b %d", len(a), len(b))
	}
	for k, v := range a {
		bvalue, ok := b[k]
		if !ok || bvalue != v {
			return fmt.Errorf("map values for key %d,  a %d, b %d", k, v, bvalue)
		}
	}
	return nil
}
