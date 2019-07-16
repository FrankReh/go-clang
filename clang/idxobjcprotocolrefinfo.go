package clang

// #include "go-clang.h"
import "C"

type IdxObjCProtocolRefInfo struct {
	c C.CXIdxObjCProtocolRefInfo
}

func (iocpri IdxObjCProtocolRefInfo) Protocol() *IdxEntityInfo {
	return newIdxEntityInfo(iocpri.c.protocol)
}

func (iocpri IdxObjCProtocolRefInfo) Cursor() Cursor {
	return Cursor{iocpri.c.cursor}
}

func (iocpri IdxObjCProtocolRefInfo) Loc() IdxLoc {
	return IdxLoc{iocpri.c.loc}
}
