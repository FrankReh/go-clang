package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import "fmt"

// Describes the exception specification of a cursor.
//
// A negative value indicates that the cursor is not a function declaration.
type ExceptionSpecification int32

const (
	// A non-function type (manually added).
	ExceptionSpecification_NonFunction ExceptionSpecification = -1
	// The cursor has no exception specification.
	ExceptionSpecification_None = C.CXCursor_ExceptionSpecificationKind_None
	// The cursor has exception specification throw()
	ExceptionSpecification_DynamicNone = C.CXCursor_ExceptionSpecificationKind_DynamicNone
	// The cursor has exception specification throw(T1, T2)
	ExceptionSpecification_Dynamic = C.CXCursor_ExceptionSpecificationKind_Dynamic
	// The cursor has exception specification throw(...).
	ExceptionSpecification_MSAny = C.CXCursor_ExceptionSpecificationKind_MSAny
	// The cursor has exception specification basic noexcept.
	ExceptionSpecification_BasicNoexcept = C.CXCursor_ExceptionSpecificationKind_BasicNoexcept
	// The cursor has exception specification computed noexcept.
	ExceptionSpecification_ComputedNoexcept = C.CXCursor_ExceptionSpecificationKind_ComputedNoexcept
	// The exception specification has not yet been evaluated.
	ExceptionSpecification_Unevaluated = C.CXCursor_ExceptionSpecificationKind_Unevaluated
	// The exception specification has not yet been instantiated.
	ExceptionSpecification_Uninstantiated = C.CXCursor_ExceptionSpecificationKind_Uninstantiated
	// The exception specification has not been parsed yet.
	ExceptionSpecification_Unparsed = C.CXCursor_ExceptionSpecificationKind_Unparsed
)

func (cesk ExceptionSpecification) Spelling() string {
	switch cesk {
	case ExceptionSpecification_NonFunction:
		return "Cursor=ExceptionSpecification_NonFunction"
	case ExceptionSpecification_None:
		return "Cursor=ExceptionSpecification_None"
	case ExceptionSpecification_DynamicNone:
		return "Cursor=ExceptionSpecification_DynamicNone"
	case ExceptionSpecification_Dynamic:
		return "Cursor=ExceptionSpecification_Dynamic"
	case ExceptionSpecification_MSAny:
		return "Cursor=ExceptionSpecification_MSAny"
	case ExceptionSpecification_BasicNoexcept:
		return "Cursor=ExceptionSpecification_BasicNoexcept"
	case ExceptionSpecification_ComputedNoexcept:
		return "Cursor=ExceptionSpecification_ComputedNoexcept"
	case ExceptionSpecification_Unevaluated:
		return "Cursor=ExceptionSpecification_Unevaluated"
	case ExceptionSpecification_Uninstantiated:
		return "Cursor=ExceptionSpecification_Uninstantiated"
	case ExceptionSpecification_Unparsed:
		return "Cursor=ExceptionSpecification_Unparsed"
	}

	return fmt.Sprintf("ExceptionSpecification unkown %d", int(cesk))
}

func (cesk ExceptionSpecification) String() string {
	return cesk.Spelling()
}
