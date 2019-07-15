package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

// Data for IndexerCallbacks#indexEntityReference.
type IdxEntityRefInfo struct {
	c C.CXIdxEntityRefInfo
}

func (ieri IdxEntityRefInfo) Kind() IdxEntityRefKind {
	return IdxEntityRefKind(ieri.c.kind)
}

// Reference cursor.
func (ieri IdxEntityRefInfo) Cursor() Cursor {
	return Cursor{ieri.c.cursor}
}

func (ieri IdxEntityRefInfo) Loc() IdxLoc {
	return IdxLoc{ieri.c.loc}
}

// The entity that gets referenced.
func (ieri IdxEntityRefInfo) ReferencedEntity() *IdxEntityInfo {
	if o := ieri.c.referencedEntity; o != nil {
		return &IdxEntityInfo{o}
	}
	return nil
}

/*
	Immediate "parent" of the reference. For example:

	\code
	Foo *var;
	\endcode

	The parent of reference of type 'Foo' is the variable 'var'.
	For references inside statement bodies of functions/methods,
	the parentEntity will be the function/method.
*/
func (ieri IdxEntityRefInfo) ParentEntity() *IdxEntityInfo {
	if o := ieri.c.parentEntity; o != nil {
		return &IdxEntityInfo{o}
	}
	return nil
}

// Lexical container context of the reference.
func (ieri IdxEntityRefInfo) Container() *IdxContainerInfo {
	if o := ieri.c.container; o != nil {
		return &IdxContainerInfo{o}
	}
	return nil
}

// Sets of symbol roles of the reference.
func (ieri IdxEntityRefInfo) Role() SymbolRole {
	return SymbolRole(ieri.c.role)
}
