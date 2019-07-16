package clang

// #include "go-clang.h"
import "C"
import (
	"reflect"
	"unsafe"
)

type IdxCXXClassDeclInfo struct {
	c C.CXIdxCXXClassDeclInfo
}

func newIdxCXXClassDeclInfo(c *C.CXIdxCXXClassDeclInfo) *IdxCXXClassDeclInfo {
	if c != nil {
		return &IdxCXXClassDeclInfo{*c}
	}
	return nil
}

func (icxxcdi IdxCXXClassDeclInfo) DeclInfo() *IdxDeclInfo {
	return newIdxDeclInfo(icxxcdi.c.declInfo)
}

func (icxxcdi IdxCXXClassDeclInfo) Bases() []*IdxBaseClassInfo {
	var s []*IdxBaseClassInfo
	gos_s := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	gos_s.Cap = int(icxxcdi.c.numBases)
	gos_s.Len = int(icxxcdi.c.numBases)
	gos_s.Data = uintptr(unsafe.Pointer(icxxcdi.c.bases))

	return s
}

func (icxxcdi IdxCXXClassDeclInfo) NumBases() uint32 {
	return uint32(icxxcdi.c.numBases)
}
