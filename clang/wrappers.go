package clang

// #include "go-clang.h"
import "C"

// Opaque pointer representing client data that will be passed through to
// various callbacks and visitors.
type ClientData struct {
	c C.CXClientData
}

// visitor callback that will receive pairs of CXCursor/CXSourceRange
type CursorAndRangeVisitor struct {
	c C.CXCursorAndRangeVisitor
}

// Uniquely identifies a CXFile, that refers to the same underlying file, across an indexing session.
type FileUniqueID struct {
	c C.CXFileUniqueID
}

// The client's data object that is associated with an AST file (PCH or module).
type IdxClientASTFile struct {
	c C.CXIdxClientASTFile
}

// The client's data object that is associated with a semantic container of entities.
type IdxClientContainer struct {
	c C.CXIdxClientContainer
}

// The client's data object that is associated with a semantic entity.
type IdxClientEntity struct {
	c C.CXIdxClientEntity
}
