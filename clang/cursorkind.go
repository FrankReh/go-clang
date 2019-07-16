package clang

// #include "go-clang.h"
import "C"
import "github.com/frankreh/go-clang/clang/cursorkind"

func IsDeclaration(ck cursorkind.Kind) bool {
	return C.clang_isDeclaration(C.enum_CXCursorKind(ck)) != 0
}

// Determine whether the given cursor kind represents a simple reference.
//
// Note that other kinds of cursors (such as expressions) can also refer to
// other cursors. Use clang_getCursorReferenced() to determine whether a
// particular cursor refers to another entity.
func IsReference(ck cursorkind.Kind) bool {
	return C.clang_isReference(C.enum_CXCursorKind(ck)) != 0
}

// Determine whether the given cursor kind represents an expression.
func IsExpression(ck cursorkind.Kind) bool {
	return C.clang_isExpression(C.enum_CXCursorKind(ck)) != 0
}

// Determine whether the given cursor kind represents a statement.
func IsStatement(ck cursorkind.Kind) bool {
	return C.clang_isStatement(C.enum_CXCursorKind(ck)) != 0
}

// Determine whether the given cursor kind represents an attribute.
func IsAttribute(ck cursorkind.Kind) bool {
	return C.clang_isAttribute(C.enum_CXCursorKind(ck)) != 0
}

// Determine whether the given cursor kind represents an invalid cursor.
func IsInvalid(ck cursorkind.Kind) bool {
	return C.clang_isInvalid(C.enum_CXCursorKind(ck)) != 0
}

// Determine whether the given cursor kind represents a translation unit.
func IsTranslationUnit(ck cursorkind.Kind) bool {
	return C.clang_isTranslationUnit(C.enum_CXCursorKind(ck)) != 0
}

// * Determine whether the given cursor represents a preprocessing element, such as a preprocessor directive or macro instantiation.
func IsPreprocessing(ck cursorkind.Kind) bool {
	return C.clang_isPreprocessing(C.enum_CXCursorKind(ck)) != 0
}

// * Determine whether the given cursor represents a currently unexposed piece of the AST (e.g., CXCursor_UnexposedStmt).
func IsUnexposed(ck cursorkind.Kind) bool {
	return C.clang_isUnexposed(C.enum_CXCursorKind(ck)) != 0
}

func CursorKindSpelling(ck cursorkind.Kind) string {
	return cx2GoString(C.clang_getCursorKindSpelling(C.enum_CXCursorKind(ck)))
}
