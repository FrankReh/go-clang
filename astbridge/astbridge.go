// Package astbridge bridges the cgo functionality of the clang package
// with the ast package structures.
//
// The ast structures do not rely on clang and can therefore be cgo free.
// The primary ast structures are also pointer free so can be serialized.
package astbridge

import (
	"errors"
	"fmt"

	"github.com/frankreh/go-clang-v5.0/ast"
	"github.com/frankreh/go-clang-v5.0/clang"
	"github.com/frankreh/go-clang-v5.0/clang/cursorkind"
)

// ClangTranslationUnit references the clang package components.
type ClangTranslationUnit struct {
	GoTu         ast.TranslationUnit
	ClangTu      *clang.TranslationUnit
	ClangTokens  []clang.Token
	ClangCursors []clang.Cursor
}

func (ctu *ClangTranslationUnit) convertClangTokens(clangTokens []clang.Token) []ast.TokenId {
	r := make([]ast.TokenId, len(clangTokens))

	for i := range r {
		tokenSpelling := ctu.ClangTu.TokenSpelling(clangTokens[i])
		token := ast.Token{
			TokenKindId: clangTokens[i].Kind(),
			TokenNameId: ctu.GoTu.TokenNameMap.Id(tokenSpelling),
		}
		tokenId := ctu.GoTu.TokenMap.Id(token)
		r[i] = tokenId
	}
	return r
}

func mapSourceLocationToIndex(tu *clang.TranslationUnit, clangTokens []clang.Token) map[clang.SourceLocation]int {
	r := make(map[clang.SourceLocation]int)

	for i := range clangTokens {
		r[tu.TokenLocation(clangTokens[i])] = i
	}
	return r
}

func (ctu *ClangTranslationUnit) Populate(tu *clang.TranslationUnit) error {
	if ctu.ClangTu != nil {
		return errors.New("Already populated")
	}

	// For some tidyness, have the "" string map to the 0 id.
	_ = ctu.GoTu.CursorNameMap.Id("")

	ctu.ClangTu = tu
	clangRootCursor := tu.TranslationUnitCursor()

	ctu.ClangTokens = tu.Tokenize(clangRootCursor.Extent())

	ctu.GoTu.TokenIds = ctu.convertClangTokens(ctu.ClangTokens)

	mapTokenIndex := mapSourceLocationToIndex(tu, ctu.ClangTokens)

	ctu.GoTu.Back = make(map[int]int)

	// Layer children to end of list, one set of children at a time.

	// Seed list with the root.
	ctu.ClangCursors = append(ctu.ClangCursors, clangRootCursor)
	ctu.GoTu.Cursors = append(ctu.GoTu.Cursors, ast.Cursor{
		CursorKindId: clangRootCursor.Kind(),
		CursorNameId: ctu.GoTu.CursorNameMap.Id(clangRootCursor.Spelling()),
		ParentIndex:  -1,
		Tokens: ast.IndexPair{ // By definition, all the clang tokens.
			Head: 0,
			Len:  len(ctu.ClangTokens),
		},

		// Index can be set manually later but we do it here
		// to better show what's going on.
		// Index: 0, // Index no longer exists in this structure.
	})

	debug := false
	if debug {
		fmt.Printf("%[1]s %[1]d\n", clang.ChildVisit_Break)
		fmt.Printf("%[1]s %[1]d\n", clang.ChildVisit_Continue)
		fmt.Printf("%[1]s %[1]d\n", clang.ChildVisit_Recurse)
	}

	// Map cursor to its index in the list being created.
	cursorsSeen := make(map[clang.Cursor]int)
	nullCursor := clang.NewNullCursor()

	// Grow the list of clang cursors by visiting the list of clang cursors.
	for parentIndex := 0; parentIndex < len(ctu.ClangCursors); parentIndex++ {

		if ctu.ClangCursors[parentIndex] == nullCursor {
			if debug {
				fmt.Printf("for loop: parentIndex %d, nullCursor\n", parentIndex)
			}
			continue
		}

		childcount := 0

		if debug {
			fmt.Printf("for loop: parentIndex %d, hash %x len %d\n",
				parentIndex,
				ctu.ClangCursors[parentIndex].HashCursor(),
				len(ctu.ClangCursors))
		}

		ctu.ClangCursors[parentIndex].Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
			childcount++

			ownIndex := len(ctu.GoTu.Cursors)

			seenIndex, seen := cursorsSeen[cursor]

			if seen {
				// If this cursor has been seen already, it means the AST clang is walking us through
				// has an additional parent for the same child. Rather than traverse this child as a new
				// cursor that needs to be recorded, along with its children, register this duplication in
				// the tree by creating a cursorkind.Back entry.

				ctu.GoTu.Back[ownIndex] = seenIndex

				backCursor := ast.Cursor{
					CursorKindId: cursorkind.Back,
					ParentIndex:  parentIndex,
				}
				// Keep the lists the same length. Use the nullCursor as a place holder.
				ctu.GoTu.Cursors = append(ctu.GoTu.Cursors, backCursor)
				ctu.ClangCursors = append(ctu.ClangCursors, nullCursor)

			} else {

				if debug {
					fmt.Printf("parentIndex %d, visiting cursor hash %x, parent hash %x\n",
						parentIndex,
						cursor.HashCursor(),
						parent.HashCursor())
					fmt.Printf("childcount %d\n", childcount)
				}

				// Start to set newCursor.Tokens.
				var tokenRange ast.IndexPair

				// Get clang tokens for this cursor long enough to find them in
				// the global tu list of tokens. It should be enough to get just
				// the first token from the cursor, but there is no libclang call
				// for that.
				if tokens := tu.Tokenize(cursor.Extent()); len(tokens) > 0 {
					// TBD work against the parent's list first to reduce the search times.
					index, ok := mapTokenIndex[tu.TokenLocation(tokens[0])]
					if !ok {
						// oken location not found in map, skipping cursor
						return clang.ChildVisit_Continue
					}
					tokenRange.Head = index
					tokenRange.Len = len(tokens)
				}
				// End to set newCursor.Tokens.

				newCursor := ast.Cursor{
					CursorKindId: cursor.Kind(),
					CursorNameId: ctu.GoTu.CursorNameMap.Id(cursor.Spelling()),
					ParentIndex:  parentIndex,

					// Index can be set manually later but we do it here
					// to better show what's going on.
					// Index:  ownIndex,
					Tokens: tokenRange,
				}

				cursorsSeen[cursor] = len(ctu.GoTu.Cursors) // Length of either list would suffice.

				// N.B. Append to the two lists back to back, else risk
				// a return getting added in between and having their
				// lengths not match.
				ctu.GoTu.Cursors = append(ctu.GoTu.Cursors, newCursor)
				ctu.ClangCursors = append(ctu.ClangCursors, cursor)
			}

			// Determining the children doesn't have to be done here.
			// It is enough that the ParentIndex was set in this visit.
			// But it is done here to better show what's going on.

			// Update parent's notion of children.
			//
			// Effectively append own index to parent's list of children.
			// First one through becomes the head, and each one through
			// increases the length by one.
			c := &ctu.GoTu.Cursors[parentIndex].Children
			if c.Head == 0 {
				// ownIndex will never by zero because list starts off with root in it,
				// so length of that list starts at 1.
				c.Head = ownIndex
			}
			c.Len++
			// End of setup that could be done later.

			return clang.ChildVisit_Continue // Continue to next sibling
		})
	}

	// These could also be called from DecodeFinish if we knew we just wanted them done
	// after gob decoding and not when struct is first populated.
	ctu.GoTu.SetBinaryOperatorNames()
	// Maybe don't call this here either.
	//ctu.GoTu.SetBackChildren()

	return nil
}
