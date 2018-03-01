package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import "fmt"

type IdxObjCContainerKind uint32

const (
	IdxObjCContainer_ForwardRef     IdxObjCContainerKind = C.CXIdxObjCContainer_ForwardRef
	IdxObjCContainer_Interface      IdxObjCContainerKind = C.CXIdxObjCContainer_Interface
	IdxObjCContainer_Implementation IdxObjCContainerKind = C.CXIdxObjCContainer_Implementation
)

func (iocck IdxObjCContainerKind) String() string {
	switch iocck {
	case IdxObjCContainer_ForwardRef:
		return "IdxObjCContainer_ForwardRef"
	case IdxObjCContainer_Interface:
		return "IdxObjCContainer_Interface"
	case IdxObjCContainer_Implementation:
		return "IdxObjCContainer_Implementation"
	}

	return fmt.Sprintf("IdxObjCContainerKind unknown %d", int(iocck))
}
