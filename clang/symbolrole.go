package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

// Roles that are attributed to symbol occurrences.
//
// Internal: this currently mirrors low 9 bits of clang::index::SymbolRole with
// higher bits zeroed. These high bits may be exposed in the future.
type SymbolRole uint32

const (
	SymbolRole_None        SymbolRole = C.CXSymbolRole_None
	SymbolRole_Declaration SymbolRole = C.CXSymbolRole_Declaration
	SymbolRole_Definition  SymbolRole = C.CXSymbolRole_Definition
	SymbolRole_Reference   SymbolRole = C.CXSymbolRole_Reference
	SymbolRole_Read        SymbolRole = C.CXSymbolRole_Read
	SymbolRole_Write       SymbolRole = C.CXSymbolRole_Write
	SymbolRole_Call        SymbolRole = C.CXSymbolRole_Call
	SymbolRole_Dynamic     SymbolRole = C.CXSymbolRole_Dynamic
	SymbolRole_AddressOf   SymbolRole = C.CXSymbolRole_AddressOf
	SymbolRole_Implicit    SymbolRole = C.CXSymbolRole_Implicit
)
