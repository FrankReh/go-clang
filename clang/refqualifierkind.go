package clang

// #include "go-clang.h"
import "C"

type RefQualifierKind uint32

const (
	// No ref-qualifier was provided.
	RefQualifier_None RefQualifierKind = C.CXRefQualifier_None

	// An lvalue ref-qualifier was provided (&).
	RefQualifier_LValue RefQualifierKind = C.CXRefQualifier_LValue

	// An rvalue ref-qualifier was provided (&&).
	RefQualifier_RValue RefQualifierKind = C.CXRefQualifier_RValue
)
