package clang

// #include "go-clang.h"
import "C"
import (
	"reflect"
	"unsafe"
)

type IdxDeclInfo struct {
	c *C.CXIdxDeclInfo
}

func newIdxDeclInfo(c *C.CXIdxDeclInfo) *IdxDeclInfo {
	if c != nil {
		return &IdxDeclInfo{c}
	}
	return nil
}

func (idi *IdxDeclInfo) ContainerDeclInfo() *IdxObjCContainerDeclInfo {
	return newIdxObjCContainerDeclInfo(C.clang_index_getObjCContainerDeclInfo(idi.c))
}

func (idi *IdxDeclInfo) InterfaceDeclInfo() *IdxObjCInterfaceDeclInfo {
	return newIdxObjCInterfaceDeclInfo(C.clang_index_getObjCInterfaceDeclInfo(idi.c))
}

func (idi *IdxDeclInfo) CategoryDeclInfo() *IdxObjCCategoryDeclInfo {
	return newIdxObjCCategoryDeclInfo(C.clang_index_getObjCCategoryDeclInfo(idi.c))
}

func (idi *IdxDeclInfo) ProtocolRefListInfo() *IdxObjCProtocolRefListInfo {
	return newIdxObjCProtocolRefListInfo(C.clang_index_getObjCProtocolRefListInfo(idi.c))
}

func (idi *IdxDeclInfo) PropertyDeclInfo() *IdxObjCPropertyDeclInfo {
	return newIdxObjCPropertyDeclInfo(C.clang_index_getObjCPropertyDeclInfo(idi.c))
}

func (idi *IdxDeclInfo) ClassDeclInfo() *IdxCXXClassDeclInfo {
	return newIdxCXXClassDeclInfo(C.clang_index_getCXXClassDeclInfo(idi.c))
}

func (idi IdxDeclInfo) EntityInfo() *IdxEntityInfo {
	return newIdxEntityInfo(idi.c.entityInfo)
}

func (idi IdxDeclInfo) Cursor() Cursor {
	return Cursor{idi.c.cursor}
}

func (idi IdxDeclInfo) Loc() IdxLoc {
	return IdxLoc{idi.c.loc}
}

func (idi IdxDeclInfo) SemanticContainer() *IdxContainerInfo {
	return newIdxContainerInfo(idi.c.semanticContainer)
}

// Generally same as #semanticContainer but can be different in cases like out-of-line C++ member functions.
func (idi IdxDeclInfo) LexicalContainer() *IdxContainerInfo {
	return newIdxContainerInfo(idi.c.lexicalContainer)
}

func (idi IdxDeclInfo) IsRedeclaration() bool {
	return idi.c.isRedeclaration != 0
}

func (idi IdxDeclInfo) IsDefinition() bool {
	return idi.c.isDefinition != 0
}

func (idi IdxDeclInfo) IsContainer() bool {
	return idi.c.isContainer != 0
}

func (idi IdxDeclInfo) DeclAsContainer() *IdxContainerInfo {
	return newIdxContainerInfo(idi.c.declAsContainer)
}

// Whether the declaration exists in code or was created implicitly by the compiler, e.g. implicit Objective-C methods for properties.
func (idi IdxDeclInfo) IsImplicit() bool {
	o := idi.c.isImplicit

	return o != C.int(0)
}

func (idi IdxDeclInfo) Attributes() []*IdxAttrInfo {
	var s []*IdxAttrInfo
	gos_s := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	gos_s.Cap = int(idi.c.numAttributes)
	gos_s.Len = int(idi.c.numAttributes)
	gos_s.Data = uintptr(unsafe.Pointer(idi.c.attributes))

	return s
}

func (idi IdxDeclInfo) NumAttributes() uint32 {
	return uint32(idi.c.numAttributes)
}

func (idi IdxDeclInfo) Flags() uint32 {
	return uint32(idi.c.flags)
}

// This type is defined by the header but not used.
// type IdxDeclInfoFlags uint32
//
// const (
// 	IdxDeclFlag_Skipped IdxDeclInfoFlags = C.CXIdxDeclFlag_Skipped
// )
