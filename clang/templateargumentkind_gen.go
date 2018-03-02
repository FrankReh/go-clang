package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

/*
	Describes the kind of a template argument.

	See the definition of llvm::clang::TemplateArgument::ArgKind for full
	element descriptions.
*/
type TemplateArgumentKind uint32

const (
	TemplateArgumentKind_Null              TemplateArgumentKind = C.CXTemplateArgumentKind_Null
	TemplateArgumentKind_Type              TemplateArgumentKind = C.CXTemplateArgumentKind_Type
	TemplateArgumentKind_Declaration       TemplateArgumentKind = C.CXTemplateArgumentKind_Declaration
	TemplateArgumentKind_NullPtr           TemplateArgumentKind = C.CXTemplateArgumentKind_NullPtr
	TemplateArgumentKind_Integral          TemplateArgumentKind = C.CXTemplateArgumentKind_Integral
	TemplateArgumentKind_Template          TemplateArgumentKind = C.CXTemplateArgumentKind_Template
	TemplateArgumentKind_TemplateExpansion TemplateArgumentKind = C.CXTemplateArgumentKind_TemplateExpansion
	TemplateArgumentKind_Expression        TemplateArgumentKind = C.CXTemplateArgumentKind_Expression
	TemplateArgumentKind_Pack              TemplateArgumentKind = C.CXTemplateArgumentKind_Pack
	TemplateArgumentKind_Invalid           TemplateArgumentKind = C.CXTemplateArgumentKind_Invalid
)
