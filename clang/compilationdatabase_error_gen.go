package clang

// #include "./clang-c/CXCompilationDatabase.h"
// #include "go-clang.h"
import "C"
import "fmt"

// Error code for Compilation Database

const CanNotLoadDatabaseErr = Error("CanNotLoadDatabase")

func convertCompilationDatabaseErrorCode(ec C.CXCompilationDatabase_Error) error {
	switch ec {
	case C.CXCompilationDatabase_NoError:
		return nil
	case C.CXCompilationDatabase_CanNotLoadDatabase:
		return CanNotLoadDatabaseErr
	}

	return fmt.Errorf("unknown CXCompilationDatabase_Error %d", int(ec))
}
