package clang

// #include "go-clang.h"
import "C"

// Data for IndexerCallbacks#indexEntityReference.
type IdxEntityRefInfo struct{ c C.CXIdxEntityRefInfo }

func (g IdxEntityRefInfo) Kind() IdxEntityRefKind {
	return IdxEntityRefKind(g.c.kind)
}

// Data for IndexerCallbacks#indexEntityReference.
type IdxEntityRefKind uint32

const (
	// The entity is referenced directly in user's code.
	IdxEntityRef_Direct IdxEntityRefKind = C.CXIdxEntityRef_Direct

	// An implicit reference, e.g. a reference of an Objective-C method via the dot syntax.
	IdxEntityRef_Implicit IdxEntityRefKind = C.CXIdxEntityRef_Implicit
)

// Reference cursor.
func (g IdxEntityRefInfo) Cursor() Cursor { return Cursor{g.c.cursor} }

// Source location used in index callbacks.
func (g IdxEntityRefInfo) Loc() IdxLoc { return IdxLoc{g.c.loc} }

// The entity that gets referenced.
func (g IdxEntityRefInfo) Entity() *IdxEntityInfo { return newIdxEntityInfo(g.c.referencedEntity) }

// Immediate parent of the reference.
func (g IdxEntityRefInfo) Parent() *IdxEntityInfo { return newIdxEntityInfo(g.c.parentEntity) }

// Lexical container context of the reference.
func (g IdxEntityRefInfo) Container() *IdxContainerInfo { return newIdxContainerInfo(g.c.container) }

// Sets of symbol roles of the reference.
func (g IdxEntityRefInfo) Role() SymbolRole { return SymbolRole(g.c.role) }
