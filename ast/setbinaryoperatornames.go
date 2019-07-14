package ast

import (
	"github.com/frankreh/go-clang/clang/cursorkind"
	"github.com/frankreh/go-clang/clang/tokenkind"
)

// SetBinaryOperatorNames sets the missing names for the
// UnaryOperator and BinaryOperator cursor types.
func (tu *TranslationUnit) SetBinaryOperatorNames() error {

	// Just visit all the nodes.
	noname, ok := tu.CursorNameMap.m[""]
	if !ok {
		panic("CursorNameMap did not include noname")
	}

	c := tu.Cursors
	for i := range c {
		if cursorkind.BinaryOperator != c[i].CursorKindId {
			continue
		}
		if noname != c[i].CursorNameId {
			continue
		}
		// Sanity check there are two children.
		if 2 != c[i].Children.Len {
			return tu.Err("BinaryOperator").Cursor(i).ChildrenCount()
		}
		child0 := c[i].Children.Head
		child1 := child0 + 1
		// Sanity check self and children have tokens.
		if 0 == c[i].Tokens.Len {
			return tu.Err("no tokens").Cursor(i)
		}
		if 0 == c[child0].Tokens.Len {
			return tu.Err("no tokens").Cursor(i).Children(child0)
		}
		if 0 == c[child1].Tokens.Len {
			return tu.Err("no tokens").Cursor(i).Children(child1)
		}
		// Sanity check the two children have single gap between tokens.
		next := c[child0].Tokens.Next()
		if 1 != c[child1].Tokens.Head-next {
			return tu.Err("BinaryOperator children should have single token gap").Cursor(i)
		}
		// Next should be the token index of the BinaryOperator.

		// Sanity check next is in range and is Punctuation and that the name is a recognized BinaryOperator.
		if next >= len(tu.TokenIds) {
			return tu.Err("TokenIds", &OutOfRangeErr{next, len(tu.TokenIds)}).Cursor(i)
		}

		tokenId := tu.TokenIds[next]
		if int(tokenId) > len(tu.TokenMap.Tokens) {
			return tu.Err("TokenMap.Tokens", &OutOfRangeErr{int(tokenId), len(tu.TokenMap.Tokens)}).Cursor(i)
		}
		token := tu.TokenMap.Tokens[tokenId]
		tokenKindId := token.TokenKindId
		tokenNameId := token.TokenNameId

		if tokenkind.Punctuation != tokenKindId {
			return tu.Err(token, "not Punctuation").Cursor(i)
		}
		if tokenNameId >= len(tu.TokenNameMap.Strings) {
			return tu.Err(token, "TokenNameMap.Strings", &OutOfRangeErr{tokenNameId, len(tu.TokenNameMap.Strings)}).Cursor(i)
		}
		tokenName := tu.TokenNameMap.Strings[tokenNameId]

		// TBD check versus valid c language binary operator names.

		cursorNameId := tu.CursorNameMap.Id(tokenName) // Adds to map if not yet included.

		// TBD may want a debugging trace statement here.
		c[i].CursorNameId = cursorNameId
	}
	return nil
}
