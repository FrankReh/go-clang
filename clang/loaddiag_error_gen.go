package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import "fmt"

const LoadDiag_UnknownErr = Error("LoadDiag_Unknown")
const LoadDiag_CannotLoadErr = Error("LoadDiag_CannotLoad")
const LoadDiag_InvalidFileErr = Error("LoadDiag_InvalidFile")

func convertLoadDiagErrorCode(ec C.enum_CXLoadDiag_Error) error {
	switch ec {
	case C.CXLoadDiag_None:
		return nil
	case C.CXLoadDiag_Unknown:
		return LoadDiag_UnknownErr
	case C.CXLoadDiag_CannotLoad:
		return LoadDiag_CannotLoadErr
	case C.CXLoadDiag_InvalidFile:
		return LoadDiag_InvalidFileErr
	}

	return fmt.Errorf("unknown CXLoadDiag_Error %d", int(ec))
}
