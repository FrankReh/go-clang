package clang

// #include "go-clang.h"
import "C"

type IdxIBOutletCollectionAttrInfo struct {
	c C.CXIdxIBOutletCollectionAttrInfo
}

func newIdxIBOutletCollectionAttrInfo(c *C.CXIdxIBOutletCollectionAttrInfo) *IdxIBOutletCollectionAttrInfo {
	if c != nil {
		return &IdxIBOutletCollectionAttrInfo{*c}
	}
	return nil
}

func (iibocai IdxIBOutletCollectionAttrInfo) AttrInfo() *IdxAttrInfo {
	return newIdxAttrInfo(iibocai.c.attrInfo)
}

func (iibocai IdxIBOutletCollectionAttrInfo) ObjcClass() *IdxEntityInfo {
	return newIdxEntityInfo(iibocai.c.objcClass)
}

func (iibocai IdxIBOutletCollectionAttrInfo) ClassCursor() Cursor {
	return Cursor{iibocai.c.classCursor}
}

func (iibocai IdxIBOutletCollectionAttrInfo) ClassLoc() IdxLoc {
	return IdxLoc{iibocai.c.classLoc}
}
