package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

// Represents the C++ access control level to a base class for a cursor with kind CX_CXXBaseSpecifier.
type AccessSpecifier uint32

const (
	AccessSpecifier_Invalid   AccessSpecifier = C.CX_CXXInvalidAccessSpecifier
	AccessSpecifier_Public    AccessSpecifier = C.CX_CXXPublic
	AccessSpecifier_Protected AccessSpecifier = C.CX_CXXProtected
	AccessSpecifier_Private   AccessSpecifier = C.CX_CXXPrivate
)
