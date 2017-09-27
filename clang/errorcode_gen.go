package clang

// #include "./clang-c/CXErrorCode.h"
// #include "go-clang.h"
import "C"
import "fmt"

// Error codes returned by libclang routines.
//
// Zero (CXError_Success) is the only error code indicating success. Other
// error codes, including not yet assigned non-zero values, indicate errors.

// Thanks Dave C., for the idea of const Errors that are based on string.
type Error string

func (e Error) Error() string { return string(e) }

// A generic error code, no further details are available.
const FailureErr = Error("Failure")

// libclang crashed while performing the requested operation.
const CrashedErr = Error("Crashed")

// The function detected that the arguments violate the function contract.
const InvalidArgumentsErr = Error("InvalidArguments")

// An AST deserialization error has occurred.
const ASTReadErr = Error("ASTRead")

// Some other error code, unexpected, was received.
const OtherErr = Error("Other")

func convertErrorCode(ec C.enum_CXErrorCode) error {
	switch ec {
	case C.CXError_Success:
		return nil
	case C.CXError_Failure:
		return FailureErr
	case C.CXError_Crashed:
		return CrashedErr
	case C.CXError_InvalidArguments:
		return InvalidArgumentsErr
	case C.CXError_ASTReadError:
		return ASTReadErr
	}

	return fmt.Errorf("unknown C.enum_CXErrorCode %d", int(ec))
}
