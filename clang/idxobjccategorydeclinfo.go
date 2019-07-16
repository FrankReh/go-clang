package clang

// #include "go-clang.h"
import "C"

type IdxObjCCategoryDeclInfo struct {
	c C.CXIdxObjCCategoryDeclInfo
}

func newIdxObjCCategoryDeclInfo(c *C.CXIdxObjCCategoryDeclInfo) *IdxObjCCategoryDeclInfo {
	if c != nil {
		return &IdxObjCCategoryDeclInfo{*c}
	}
	return nil
}

func (ioccdi IdxObjCCategoryDeclInfo) ContainerInfo() *IdxObjCContainerDeclInfo {
	return newIdxObjCContainerDeclInfo(ioccdi.c.containerInfo)
}

func (ioccdi IdxObjCCategoryDeclInfo) ObjcClass() *IdxEntityInfo {
	return newIdxEntityInfo(ioccdi.c.objcClass)
}

func (ioccdi IdxObjCCategoryDeclInfo) ClassCursor() Cursor {
	return Cursor{ioccdi.c.classCursor}
}

func (ioccdi IdxObjCCategoryDeclInfo) ClassLoc() IdxLoc {
	return IdxLoc{ioccdi.c.classLoc}
}

func (ioccdi IdxObjCCategoryDeclInfo) Protocols() *IdxObjCProtocolRefListInfo {
	return newIdxObjCProtocolRefListInfo(ioccdi.c.protocols)
}
