package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import "fmt"

type IdxEntityKind uint32

const (
	IdxEntity_Unexposed             IdxEntityKind = C.CXIdxEntity_Unexposed
	IdxEntity_Typedef               IdxEntityKind = C.CXIdxEntity_Typedef
	IdxEntity_Function              IdxEntityKind = C.CXIdxEntity_Function
	IdxEntity_Variable              IdxEntityKind = C.CXIdxEntity_Variable
	IdxEntity_Field                 IdxEntityKind = C.CXIdxEntity_Field
	IdxEntity_EnumConstant          IdxEntityKind = C.CXIdxEntity_EnumConstant
	IdxEntity_ObjCClass             IdxEntityKind = C.CXIdxEntity_ObjCClass
	IdxEntity_ObjCProtocol          IdxEntityKind = C.CXIdxEntity_ObjCProtocol
	IdxEntity_ObjCCategory          IdxEntityKind = C.CXIdxEntity_ObjCCategory
	IdxEntity_ObjCInstanceMethod    IdxEntityKind = C.CXIdxEntity_ObjCInstanceMethod
	IdxEntity_ObjCClassMethod       IdxEntityKind = C.CXIdxEntity_ObjCClassMethod
	IdxEntity_ObjCProperty          IdxEntityKind = C.CXIdxEntity_ObjCProperty
	IdxEntity_ObjCIvar              IdxEntityKind = C.CXIdxEntity_ObjCIvar
	IdxEntity_Enum                  IdxEntityKind = C.CXIdxEntity_Enum
	IdxEntity_Struct                IdxEntityKind = C.CXIdxEntity_Struct
	IdxEntity_Union                 IdxEntityKind = C.CXIdxEntity_Union
	IdxEntity_CXXClass              IdxEntityKind = C.CXIdxEntity_CXXClass
	IdxEntity_CXXNamespace          IdxEntityKind = C.CXIdxEntity_CXXNamespace
	IdxEntity_CXXNamespaceAlias     IdxEntityKind = C.CXIdxEntity_CXXNamespaceAlias
	IdxEntity_CXXStaticVariable     IdxEntityKind = C.CXIdxEntity_CXXStaticVariable
	IdxEntity_CXXStaticMethod       IdxEntityKind = C.CXIdxEntity_CXXStaticMethod
	IdxEntity_CXXInstanceMethod     IdxEntityKind = C.CXIdxEntity_CXXInstanceMethod
	IdxEntity_CXXConstructor        IdxEntityKind = C.CXIdxEntity_CXXConstructor
	IdxEntity_CXXDestructor         IdxEntityKind = C.CXIdxEntity_CXXDestructor
	IdxEntity_CXXConversionFunction IdxEntityKind = C.CXIdxEntity_CXXConversionFunction
	IdxEntity_CXXTypeAlias          IdxEntityKind = C.CXIdxEntity_CXXTypeAlias
	IdxEntity_CXXInterface          IdxEntityKind = C.CXIdxEntity_CXXInterface
)

func (iek IdxEntityKind) IsEntityObjCContainerKind() bool {
	o := C.clang_index_isEntityObjCContainerKind(C.CXIdxEntityKind(iek))

	return o != C.int(0)
}

func (iek IdxEntityKind) String() string {
	switch iek {
	case IdxEntity_Unexposed:
		return "IdxEntity_Unexposed"
	case IdxEntity_Typedef:
		return "IdxEntity_Typedef"
	case IdxEntity_Function:
		return "IdxEntity_Function"
	case IdxEntity_Variable:
		return "IdxEntity_Variable"
	case IdxEntity_Field:
		return "IdxEntity_Field"
	case IdxEntity_EnumConstant:
		return "IdxEntity_EnumConstant"
	case IdxEntity_ObjCClass:
		return "IdxEntity_ObjCClass"
	case IdxEntity_ObjCProtocol:
		return "IdxEntity_ObjCProtocol"
	case IdxEntity_ObjCCategory:
		return "IdxEntity_ObjCCategory"
	case IdxEntity_ObjCInstanceMethod:
		return "IdxEntity_ObjCInstanceMethod"
	case IdxEntity_ObjCClassMethod:
		return "IdxEntity_ObjCClassMethod"
	case IdxEntity_ObjCProperty:
		return "IdxEntity_ObjCProperty"
	case IdxEntity_ObjCIvar:
		return "IdxEntity_ObjCIvar"
	case IdxEntity_Enum:
		return "IdxEntity_Enum"
	case IdxEntity_Struct:
		return "IdxEntity_Struct"
	case IdxEntity_Union:
		return "IdxEntity_Union"
	case IdxEntity_CXXClass:
		return "IdxEntity_CXXClass"
	case IdxEntity_CXXNamespace:
		return "IdxEntity_CXXNamespace"
	case IdxEntity_CXXNamespaceAlias:
		return "IdxEntity_CXXNamespaceAlias"
	case IdxEntity_CXXStaticVariable:
		return "IdxEntity_CXXStaticVariable"
	case IdxEntity_CXXStaticMethod:
		return "IdxEntity_CXXStaticMethod"
	case IdxEntity_CXXInstanceMethod:
		return "IdxEntity_CXXInstanceMethod"
	case IdxEntity_CXXConstructor:
		return "IdxEntity_CXXConstructor"
	case IdxEntity_CXXDestructor:
		return "IdxEntity_CXXDestructor"
	case IdxEntity_CXXConversionFunction:
		return "IdxEntity_CXXConversionFunction"
	case IdxEntity_CXXTypeAlias:
		return "IdxEntity_CXXTypeAlias"
	case IdxEntity_CXXInterface:
		return "IdxEntity_CXXInterface"
	}

	return fmt.Sprintf("IdxEntityKind unknown %d", int(iek))
}
