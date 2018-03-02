package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

// Describes the exception specification of a cursor.
//
// A negative value indicates that the cursor is not a function declaration.
type ExceptionSpecification int32

const (
	// A non-function type (manually added).
	ExceptionSpecification_NonFunction ExceptionSpecification = -1
	// The cursor has no exception specification.
	ExceptionSpecification_None ExceptionSpecification = C.CXCursor_ExceptionSpecificationKind_None
	// The cursor has exception specification throw()
	ExceptionSpecification_DynamicNone ExceptionSpecification = C.CXCursor_ExceptionSpecificationKind_DynamicNone
	// The cursor has exception specification throw(T1, T2)
	ExceptionSpecification_Dynamic ExceptionSpecification = C.CXCursor_ExceptionSpecificationKind_Dynamic
	// The cursor has exception specification throw(...).
	ExceptionSpecification_MSAny ExceptionSpecification = C.CXCursor_ExceptionSpecificationKind_MSAny
	// The cursor has exception specification basic noexcept.
	ExceptionSpecification_BasicNoexcept ExceptionSpecification = C.CXCursor_ExceptionSpecificationKind_BasicNoexcept
	// The cursor has exception specification computed noexcept.
	ExceptionSpecification_ComputedNoexcept ExceptionSpecification = C.CXCursor_ExceptionSpecificationKind_ComputedNoexcept
	// The exception specification has not yet been evaluated.
	ExceptionSpecification_Unevaluated ExceptionSpecification = C.CXCursor_ExceptionSpecificationKind_Unevaluated
	// The exception specification has not yet been instantiated.
	ExceptionSpecification_Uninstantiated ExceptionSpecification = C.CXCursor_ExceptionSpecificationKind_Uninstantiated
	// The exception specification has not been parsed yet.
	ExceptionSpecification_Unparsed ExceptionSpecification = C.CXCursor_ExceptionSpecificationKind_Unparsed
)
