package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import "fmt"

/*
	List the possible error codes for clang_Type_getSizeOf,
	clang_Type_getAlignOf, clang_Type_getOffsetOf and
	clang_Cursor_getOffsetOf.

	A value of this enumeration type can be returned if the target type is not
	a valid argument to sizeof, alignof or offsetof.
*/

const TypeLayout_InvalidErr = Error("InvalidTypeLayout")
const TypeLayout_IncompleteErr = Error("IncompleteTypeLayout")
const TypeLayout_DependentErr = Error("DependentTypeLayout")
const TypeLayout_NotConstantSizeErr = Error("NotConstantSizeTypeLayout")
const TypeLayout_InvalidFieldNameErr = Error("InvalidFieldNameTypeLayout")

func convertTypeLayoutError(r C.longlong) (uint64, error) {
	if r >= 0 {
		return uint64(r), nil
	}
	switch r {
	case C.CXTypeLayoutError_Invalid:
		return 0, TypeLayout_InvalidErr
	case C.CXTypeLayoutError_Incomplete:
		return 0, TypeLayout_IncompleteErr
	case C.CXTypeLayoutError_Dependent:
		return 0, TypeLayout_DependentErr
	case C.CXTypeLayoutError_NotConstantSize:
		return 0, TypeLayout_NotConstantSizeErr
	case C.CXTypeLayoutError_InvalidFieldName:
		return 0, TypeLayout_InvalidFieldNameErr
	}

	return 0, fmt.Errorf("unknown CXTypeLayoutError %d", r)
}
