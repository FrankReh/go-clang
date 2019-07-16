package clang

// #include "./clang-c/CXCompilationDatabase.h"
// #include "go-clang.h"
import "C"
import (
	"fmt"
	"unsafe"
)

// A compilation database holds all information used to compile files in a
// project. For each file in the database, it can be queried for the working
// directory or the command line used for the compiler invocation.
//
// Must be freed by clang_CompilationDatabase_dispose
type CompilationDatabase struct {
	c C.CXCompilationDatabase
}

/*
	Creates a compilation database from the database found in directory
	buildDir. For example, CMake can output a compile_commands.json which can
	be used to build the database.

	It must be freed by clang_CompilationDatabase_dispose.
*/
func FromDirectory(buildDir string) (CompilationDatabase, error) {
	var errorCode C.CXCompilationDatabase_Error

	c_buildDir := C.CString(buildDir)
	defer C.free(unsafe.Pointer(c_buildDir))

	o := CompilationDatabase{C.clang_CompilationDatabase_fromDirectory(c_buildDir, &errorCode)}

	return o, convertCompilationDatabaseErrorCode(errorCode)
}

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

// Free the given compilation database
func (cd CompilationDatabase) Dispose() {
	C.clang_CompilationDatabase_dispose(cd.c)
}

// Find the compile commands used for a file.
func (cd CompilationDatabase) CompileCommands(completeFileName string) []CompileCommand {
	c_completeFileName := C.CString(completeFileName)
	defer C.free(unsafe.Pointer(c_completeFileName))

	return convertCompileCommandsAndDispose(C.clang_CompilationDatabase_getCompileCommands(cd.c, c_completeFileName))
}

// Get all the compile commands in the given compilation database.
func (cd CompilationDatabase) AllCompileCommands() []CompileCommand {
	return convertCompileCommandsAndDispose(C.clang_CompilationDatabase_getAllCompileCommands(cd.c))
}
