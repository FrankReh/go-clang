package clang

// #include "go-clang.h"
import "C"

type IdxObjCContainerDeclInfo struct {
	c C.CXIdxObjCContainerDeclInfo
}

func newIdxObjCContainerDeclInfo(c *C.CXIdxObjCContainerDeclInfo) *IdxObjCContainerDeclInfo {
	if c != nil {
		return &IdxObjCContainerDeclInfo{*c}
	}
	return nil
}

func (ioccdi IdxObjCContainerDeclInfo) DeclInfo() *IdxDeclInfo {
	return newIdxDeclInfo(ioccdi.c.declInfo)
}

func (ioccdi IdxObjCContainerDeclInfo) Kind() IdxObjCContainerKind {
	return IdxObjCContainerKind(ioccdi.c.kind)
}

type IdxObjCContainerKind uint32

const (
	IdxObjCContainer_ForwardRef     IdxObjCContainerKind = C.CXIdxObjCContainer_ForwardRef
	IdxObjCContainer_Interface      IdxObjCContainerKind = C.CXIdxObjCContainer_Interface
	IdxObjCContainer_Implementation IdxObjCContainerKind = C.CXIdxObjCContainer_Implementation
)
