package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

// Data for IndexerCallbacks#indexEntityReference.
type IdxEntityRefKind uint32

const (
	// The entity is referenced directly in user's code.
	IdxEntityRef_Direct IdxEntityRefKind = C.CXIdxEntityRef_Direct

	// An implicit reference, e.g. a reference of an Objective-C method via the dot syntax.
	IdxEntityRef_Implicit IdxEntityRefKind = C.CXIdxEntityRef_Implicit
)
