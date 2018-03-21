package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import "github.com/frankreh/go-clang-v5.0/clang/tokenkind"

// Describes a single preprocessing token.
type Token struct {
	c C.CXToken
}

// Determine the kind of the given token.
func (t Token) Kind() tokenkind.Kind {
	return tokenkind.MustValidate(int(C.clang_getTokenKind(t.c)))
}
