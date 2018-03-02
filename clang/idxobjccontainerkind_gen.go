package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

type IdxObjCContainerKind uint32

const (
	IdxObjCContainer_ForwardRef     IdxObjCContainerKind = C.CXIdxObjCContainer_ForwardRef
	IdxObjCContainer_Interface      IdxObjCContainerKind = C.CXIdxObjCContainer_Interface
	IdxObjCContainer_Implementation IdxObjCContainerKind = C.CXIdxObjCContainer_Implementation
)
