package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import (
	"fmt"
	"strings"
)

// Flags that can be passed to clang_codeCompleteAt() to modify its behavior.
//
// The enumerators in this enumeration can be bitwise-OR'd together to provide
// multiple options to clang_codeCompleteAt().
type CodeComplete_Flags uint32

const (
	// Whether to include macros within the set of code completions returned.
	CodeComplete_IncludeMacros CodeComplete_Flags = C.CXCodeComplete_IncludeMacros

	// Whether to include code patterns for language constructs within the set of code completions, e.g., for loops.
	CodeComplete_IncludeCodePatterns CodeComplete_Flags = C.CXCodeComplete_IncludeCodePatterns

	// Whether to include brief documentation within the set of code completions returned.
	CodeComplete_IncludeBriefComments CodeComplete_Flags = C.CXCodeComplete_IncludeBriefComments
)

func (ccf CodeComplete_Flags) String() string {
	var r []string
	for _, t := range []struct {
		flag CodeComplete_Flags
		name string
	}{
		{CodeComplete_IncludeMacros, "IncludeMacros"},
		{CodeComplete_IncludeCodePatterns, "IncludeCodePatterns"},
		{CodeComplete_IncludeBriefComments, "IncludeBriefComments"},
	} {
		if ccf&t.flag == 0 {
			continue
		}
		ccf &^= t.flag
		r = append(r, t.name)
	}
	if ccf != 0 {
		// This cast to a large intrinsic is important; it avoids recursive calls to String().
		r = append(r, fmt.Sprintf("additional-bits(%x)", uint64(ccf)))
	}
	return strings.Join(r, ",")
}
