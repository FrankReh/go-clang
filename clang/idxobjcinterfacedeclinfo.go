package clang

// #include "go-clang.h"
import "C"

type IdxObjCInterfaceDeclInfo struct {
	c C.CXIdxObjCInterfaceDeclInfo
}

func newIdxObjCInterfaceDeclInfo(c *C.CXIdxObjCInterfaceDeclInfo) *IdxObjCInterfaceDeclInfo {
	if c != nil {
		return &IdxObjCInterfaceDeclInfo{*c}
	}
	return nil
}

func (iocidi IdxObjCInterfaceDeclInfo) ContainerInfo() *IdxObjCContainerDeclInfo {
	return newIdxObjCContainerDeclInfo(iocidi.c.containerInfo)
}

func (iocidi IdxObjCInterfaceDeclInfo) SuperInfo() *IdxBaseClassInfo {
	return newIdxBaseClassInfo(iocidi.c.superInfo)
}

func (iocidi IdxObjCInterfaceDeclInfo) Protocols() *IdxObjCProtocolRefListInfo {
	return newIdxObjCProtocolRefListInfo(iocidi.c.protocols)
}
