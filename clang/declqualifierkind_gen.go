package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

// 'Qualifiers' written next to the return and parameter types in Objective-C method declarations.
type DeclQualifierKind uint32

const (
	// This is meant to be a bitmask type. It's a mistake to have a value for 0.
	// So manually comment it out.
	//DeclQualifier_None   DeclQualifierKind = C.CXObjCDeclQualifier_None
	DeclQualifier_In     DeclQualifierKind = C.CXObjCDeclQualifier_In
	DeclQualifier_Inout  DeclQualifierKind = C.CXObjCDeclQualifier_Inout
	DeclQualifier_Out    DeclQualifierKind = C.CXObjCDeclQualifier_Out
	DeclQualifier_Bycopy DeclQualifierKind = C.CXObjCDeclQualifier_Bycopy
	DeclQualifier_Byref  DeclQualifierKind = C.CXObjCDeclQualifier_Byref
	DeclQualifier_Oneway DeclQualifierKind = C.CXObjCDeclQualifier_Oneway
)
