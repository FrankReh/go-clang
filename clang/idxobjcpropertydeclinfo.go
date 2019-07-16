package clang

// #include "go-clang.h"
import "C"

type IdxObjCPropertyDeclInfo struct {
	c C.CXIdxObjCPropertyDeclInfo
}

func newIdxObjCPropertyDeclInfo(c *C.CXIdxObjCPropertyDeclInfo) *IdxObjCPropertyDeclInfo {
	if c != nil {
		return &IdxObjCPropertyDeclInfo{*c}
	}
	return nil
}

func (iocpdi IdxObjCPropertyDeclInfo) DeclInfo() *IdxDeclInfo {
	return newIdxDeclInfo(iocpdi.c.declInfo)
}

func (iocpdi IdxObjCPropertyDeclInfo) Getter() *IdxEntityInfo {
	return newIdxEntityInfo(iocpdi.c.getter)
}

func (iocpdi IdxObjCPropertyDeclInfo) Setter() *IdxEntityInfo {
	return newIdxEntityInfo(iocpdi.c.setter)
}
