package clang

// #include "go-clang.h"
import "C"

type IdxBaseClassInfo struct {
	c C.CXIdxBaseClassInfo
}

func newIdxBaseClassInfo(c *C.CXIdxBaseClassInfo) *IdxBaseClassInfo {
	if c != nil {
		return &IdxBaseClassInfo{*c}
	}
	return nil
}

func (ibci IdxBaseClassInfo) Base() *IdxEntityInfo {
	return newIdxEntityInfo(ibci.c.base)
}

func (ibci IdxBaseClassInfo) Cursor() Cursor {
	return Cursor{ibci.c.cursor}
}

func (ibci IdxBaseClassInfo) Loc() IdxLoc {
	return IdxLoc{ibci.c.loc}
}
