package clang

// #include "go-clang.h"
import "C"

type IdxAttrInfo struct {
	c *C.CXIdxAttrInfo
}

func newIdxAttrInfo(c *C.CXIdxAttrInfo) *IdxAttrInfo {
	if c != nil {
		return &IdxAttrInfo{c}
	}
	return nil
}

func (iai *IdxAttrInfo) IBOutletCollectionAttrInfo() *IdxIBOutletCollectionAttrInfo {
	return newIdxIBOutletCollectionAttrInfo(C.clang_index_getIBOutletCollectionAttrInfo(iai.c))
}

func (iai IdxAttrInfo) Kind() IdxAttrKind {
	return IdxAttrKind(iai.c.kind)
}

func (iai IdxAttrInfo) Cursor() Cursor {
	return Cursor{iai.c.cursor}
}

func (iai IdxAttrInfo) Loc() IdxLoc {
	return IdxLoc{iai.c.loc}
}

type IdxAttrKind uint32

const (
	IdxAttr_Unexposed          IdxAttrKind = C.CXIdxAttr_Unexposed
	IdxAttr_IBAction           IdxAttrKind = C.CXIdxAttr_IBAction
	IdxAttr_IBOutlet           IdxAttrKind = C.CXIdxAttr_IBOutlet
	IdxAttr_IBOutletCollection IdxAttrKind = C.CXIdxAttr_IBOutletCollection
)
