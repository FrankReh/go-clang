package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import "fmt"

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

func (iecxxtk IdxEntityCXXTemplateKind) String() string {
	switch iecxxtk {
	case IdxEntity_NonTemplate:
		return "IdxEntity_NonTemplate"
	case IdxEntity_Template:
		return "IdxEntity_Template"
	case IdxEntity_TemplatePartialSpecialization:
		return "IdxEntity_TemplatePartialSpecialization"
	case IdxEntity_TemplateSpecialization:
		return "IdxEntity_TemplateSpecialization"
	}

	return fmt.Sprintf("IdxEntityCXXTemplateKind unknown %d", int(iecxxtk))
}
