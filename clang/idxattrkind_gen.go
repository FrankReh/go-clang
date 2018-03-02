package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

type IdxAttrKind uint32

const (
	IdxAttr_Unexposed          IdxAttrKind = C.CXIdxAttr_Unexposed
	IdxAttr_IBAction           IdxAttrKind = C.CXIdxAttr_IBAction
	IdxAttr_IBOutlet           IdxAttrKind = C.CXIdxAttr_IBOutlet
	IdxAttr_IBOutletCollection IdxAttrKind = C.CXIdxAttr_IBOutletCollection
)
