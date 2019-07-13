package ast

import (
	"fmt"

	"github.com/frankreh/go-clang-v5.0/clang/tokenkind"
)

// Token is the pure Go version of the clang Token.
type Token struct {
	TokenKindId tokenkind.Kind
	TokenNameId int // Id into TokenNameMap
}

type TokenId int // TokenMap allows Token <-> TokenId

// TokenMap, same structure as StringMap
type TokenMap struct {
	m      map[Token]TokenId
	Tokens []Token
}

/*
func (tm TokenMap) String() string {
	return fmt.Sprintf("%v", tm.Tokens)
}
*/
func (tm TokenMap) GoString() string {
	return fmt.Sprintf("{Tokens: %#v}", tm.Tokens)
}

// Len returns number of Tokens mapped. Valid indexes will be 0..Len-1.
func (tm *TokenMap) Len() int {
	return len(tm.Tokens)
}

func (tm *TokenMap) DecodeFinish() {
	if tm.m == nil {
		tm.m = make(map[Token]TokenId, len(tm.Tokens))
	}
	for i, t := range tm.Tokens {
		tm.m[t] = TokenId(i) // cast index to TokenId
	}
}

func (a *TokenMap) AssertEqual(b *TokenMap) error {
	if a == b {
		return nil
	}
	if err := a.assertEqualSlice(b); err != nil {
		return err
	}
	if err := a.assertEqualMap(b); err != nil {
		return err
	}
	return nil
}
func (a *TokenMap) assertEqualSlice(b *TokenMap) error {
	if len(a.Tokens) != len(b.Tokens) {
		return fmt.Errorf("TokenMap unequal slice lengths, %d %d",
			len(a.Tokens), len(b.Tokens))
	}
	for i, v := range a.Tokens {
		v2 := b.Tokens[i]
		if v != v2 {
			return fmt.Errorf("TokenMap unequal slice entry, %d %#v %#v",
				i, v, v2)
		}
	}
	return nil
}
func (a *TokenMap) assertEqualMap(b *TokenMap) error {
	if a.m == nil && b.m == nil {
		return nil
	}
	if len(a.m) != len(b.m) {
		return fmt.Errorf("TokenMap unequal map lengths, %d %d",
			len(a.m), len(b.m))
	}
	for k, v := range a.m {
		v2 := b.m[k]
		if v != v2 {
			return fmt.Errorf("TokenMap unequal map entry, %d %#v %#v",
				k, v, v2)
		}
	}
	return nil
}

func (tm *TokenMap) Init(Tokens []Token) {
	tm.m = make(map[Token]TokenId)
	tm.Tokens = make([]Token, len(Tokens))
	for i, t := range Tokens {
		tm.m[t] = TokenId(i) // cast index to TokenId
		tm.Tokens[i] = t
	}
}

// Id returns the id for this Token.
// First id will be zero.
func (tm *TokenMap) Id(t Token) TokenId {
	if tm.m == nil {
		tm.m = make(map[Token]TokenId)
	}
	id, ok := tm.m[t]
	if ok {
		return id
	}
	id = TokenId(len(tm.Tokens)) // cast len() int to next TokenId
	tm.m[t] = id
	tm.Tokens = append(tm.Tokens, t)
	return id
}

// String returns the string for the given id.
// Panic if id is out of range.
func (tm *TokenMap) ToToken(id TokenId) Token {
	return tm.Tokens[id]
}
