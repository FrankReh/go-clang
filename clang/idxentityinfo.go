package clang

// #include "go-clang.h"
import "C"
import (
	"reflect"
	"unsafe"
)

type IdxEntityInfo struct {
	c *C.CXIdxEntityInfo
}

func newIdxEntityInfo(c *C.CXIdxEntityInfo) *IdxEntityInfo {
	if c != nil {
		return &IdxEntityInfo{c}
	}
	return nil
}

// For retrieving a custom CXIdxClientEntity attached to an entity.
func (iei *IdxEntityInfo) ClientEntity() IdxClientEntity {
	return IdxClientEntity{C.clang_index_getClientEntity(iei.c)}
}

// For setting a custom CXIdxClientEntity attached to an entity.
func (iei *IdxEntityInfo) SetClientEntity(ice IdxClientEntity) {
	C.clang_index_setClientEntity(iei.c, ice.c)
}

func (iei IdxEntityInfo) Kind() IdxEntityKind {
	return IdxEntityKind(iei.c.kind)
}

func (iei IdxEntityInfo) TemplateKind() IdxEntityCXXTemplateKind {
	return IdxEntityCXXTemplateKind(iei.c.templateKind)
}

// Extra C++ template information for an entity. This can apply to:
// CXIdxEntity_Function CXIdxEntity_CXXClass CXIdxEntity_CXXStaticMethod
// CXIdxEntity_CXXInstanceMethod CXIdxEntity_CXXConstructor
// CXIdxEntity_CXXConversionFunction CXIdxEntity_CXXTypeAlias
type IdxEntityCXXTemplateKind uint32

const (
	IdxEntity_NonTemplate                   IdxEntityCXXTemplateKind = C.CXIdxEntity_NonTemplate
	IdxEntity_Template                      IdxEntityCXXTemplateKind = C.CXIdxEntity_Template
	IdxEntity_TemplatePartialSpecialization IdxEntityCXXTemplateKind = C.CXIdxEntity_TemplatePartialSpecialization
	IdxEntity_TemplateSpecialization        IdxEntityCXXTemplateKind = C.CXIdxEntity_TemplateSpecialization
)

func (iei IdxEntityInfo) Lang() IdxEntityLanguage {
	return IdxEntityLanguage(iei.c.lang)
}

type IdxEntityLanguage uint32

const (
	IdxEntityLang_None  IdxEntityLanguage = C.CXIdxEntityLang_None
	IdxEntityLang_C     IdxEntityLanguage = C.CXIdxEntityLang_C
	IdxEntityLang_ObjC  IdxEntityLanguage = C.CXIdxEntityLang_ObjC
	IdxEntityLang_CXX   IdxEntityLanguage = C.CXIdxEntityLang_CXX
	IdxEntityLang_Swift IdxEntityLanguage = C.CXIdxEntityLang_Swift
)

func (iei IdxEntityInfo) Name() string {
	return C.GoString(iei.c.name)
}

func (iei IdxEntityInfo) USR() string {
	return C.GoString(iei.c.USR)
}

func (iei IdxEntityInfo) Cursor() Cursor {
	return Cursor{iei.c.cursor}
}

func (iei IdxEntityInfo) Attributes() []*IdxAttrInfo {
	var s []*IdxAttrInfo
	gos_s := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	gos_s.Cap = int(iei.c.numAttributes)
	gos_s.Len = int(iei.c.numAttributes)
	gos_s.Data = uintptr(unsafe.Pointer(iei.c.attributes))

	return s
}

func (iei IdxEntityInfo) NumAttributes() uint32 {
	return uint32(iei.c.numAttributes)
}
