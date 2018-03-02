package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

type VisibilityKind uint32

const (
	// This value indicates that no visibility information is available for a provided CXCursor.
	Visibility_Invalid VisibilityKind = C.CXVisibility_Invalid

	// Symbol not seen by the linker.
	Visibility_Hidden VisibilityKind = C.CXVisibility_Hidden

	// Symbol seen by the linker but resolves to a symbol inside this object.
	Visibility_Protected VisibilityKind = C.CXVisibility_Protected

	// Symbol seen by the linker and acts like a normal symbol.
	Visibility_Default VisibilityKind = C.CXVisibility_Default
)
