package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

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
