package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import "fmt"

type IdxAttrKind uint32

const (
	IdxAttr_Unexposed          IdxAttrKind = C.CXIdxAttr_Unexposed
	IdxAttr_IBAction           IdxAttrKind = C.CXIdxAttr_IBAction
	IdxAttr_IBOutlet           IdxAttrKind = C.CXIdxAttr_IBOutlet
	IdxAttr_IBOutletCollection IdxAttrKind = C.CXIdxAttr_IBOutletCollection
)

func (iak IdxAttrKind) String() string {
	switch iak {
	case IdxAttr_Unexposed:
		return "IdxAttr_Unexposed"
	case IdxAttr_IBAction:
		return "IdxAttr_IBAction"
	case IdxAttr_IBOutlet:
		return "IdxAttr_IBOutlet"
	case IdxAttr_IBOutletCollection:
		return "IdxAttr_IBOutletCollection"
	}

	return fmt.Sprintf("IdxAttrKind unknown %d", int(iak))
}
