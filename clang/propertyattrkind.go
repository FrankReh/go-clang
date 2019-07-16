package clang

// #include "go-clang.h"
import "C"

// Property attributes for a CXCursor_ObjCPropertyDecl.
type PropertyAttrKind uint32

const (
	PropertyAttr_noattr            PropertyAttrKind = C.CXObjCPropertyAttr_noattr
	PropertyAttr_readonly          PropertyAttrKind = C.CXObjCPropertyAttr_readonly
	PropertyAttr_getter            PropertyAttrKind = C.CXObjCPropertyAttr_getter
	PropertyAttr_assign            PropertyAttrKind = C.CXObjCPropertyAttr_assign
	PropertyAttr_readwrite         PropertyAttrKind = C.CXObjCPropertyAttr_readwrite
	PropertyAttr_retain            PropertyAttrKind = C.CXObjCPropertyAttr_retain
	PropertyAttr_copy              PropertyAttrKind = C.CXObjCPropertyAttr_copy
	PropertyAttr_nonatomic         PropertyAttrKind = C.CXObjCPropertyAttr_nonatomic
	PropertyAttr_setter            PropertyAttrKind = C.CXObjCPropertyAttr_setter
	PropertyAttr_atomic            PropertyAttrKind = C.CXObjCPropertyAttr_atomic
	PropertyAttr_weak              PropertyAttrKind = C.CXObjCPropertyAttr_weak
	PropertyAttr_strong            PropertyAttrKind = C.CXObjCPropertyAttr_strong
	PropertyAttr_unsafe_unretained PropertyAttrKind = C.CXObjCPropertyAttr_unsafe_unretained
	PropertyAttr_class             PropertyAttrKind = C.CXObjCPropertyAttr_class
)
