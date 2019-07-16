package clang

// #include "go-clang.h"
import "C"

// Describe the linkage of the entity referred to by a cursor.
type LinkageKind uint32

const (
	// This value indicates that no linkage information is available for a provided CXCursor.
	Linkage_Invalid LinkageKind = C.CXLinkage_Invalid

	// This is the linkage for variables, parameters, and so on that have automatic storage. This covers normal (non-extern) local variables.
	Linkage_NoLinkage LinkageKind = C.CXLinkage_NoLinkage

	// This is the linkage for static variables and static functions.
	Linkage_Internal LinkageKind = C.CXLinkage_Internal

	// This is the linkage for entities with external linkage that live in C++ anonymous namespaces.
	Linkage_UniqueExternal LinkageKind = C.CXLinkage_UniqueExternal

	// This is the linkage for entities with true, external linkage.
	Linkage_External LinkageKind = C.CXLinkage_External
)
