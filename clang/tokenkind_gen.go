package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

// Describes a kind of token.
type TokenKind uint32

const (
	// A token that contains some kind of punctuation.
	Token_Punctuation TokenKind = C.CXToken_Punctuation

	// A language keyword.
	Token_Keyword TokenKind = C.CXToken_Keyword

	// An identifier (that is not a keyword).
	Token_Identifier TokenKind = C.CXToken_Identifier

	// A numeric, string, or character literal.
	Token_Literal TokenKind = C.CXToken_Literal

	// A comment.
	Token_Comment TokenKind = C.CXToken_Comment
)
