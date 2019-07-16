package clang

// #include "go-clang.h"
import "C"

type IdxContainerInfo struct {
	c *C.CXIdxContainerInfo
}

func newIdxContainerInfo(c *C.CXIdxContainerInfo) *IdxContainerInfo {
	if c != nil {
		return &IdxContainerInfo{c}
	}
	return nil
}

// For retrieving a custom CXIdxClientContainer attached to a container.
func (ici *IdxContainerInfo) ClientContainer() IdxClientContainer {
	return IdxClientContainer{C.clang_index_getClientContainer(ici.c)}
}

// For setting a custom CXIdxClientContainer attached to a container.
func (ici *IdxContainerInfo) SetClientContainer(icc IdxClientContainer) {
	C.clang_index_setClientContainer(ici.c, icc.c)
}

func (ici IdxContainerInfo) Cursor() Cursor {
	return Cursor{ici.c.cursor}
}
