package clang

// #include "./clang-c/CXCompilationDatabase.h"
// #include "go-clang.h"
import "C"

/*
	Contains the results of a search in the compilation database

	When searching for the compile command for a file, the compilation db can
	return several commands, as the file may have been compiled with
	different options in different places of the project. This choice of compile
	commands is wrapped in this opaque data structure. It must be freed by
	clang_CompileCommands_dispose.
*/

func convertCompileCommandsAndDispose(c C.CXCompileCommands) []CompileCommand {
	n := int(C.clang_CompileCommands_getSize(c))

	r := make([]CompileCommand, n)

	for i := range r {
		r[i] = newCompileCommand(C.clang_CompileCommands_getCommand(c, C.uint(i)))
	}

	C.clang_CompileCommands_dispose(c)
	return r
}
