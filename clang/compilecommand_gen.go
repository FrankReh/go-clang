package clang

// #include "./clang-c/CXCompilationDatabase.h"
// #include "go-clang.h"
import "C"

type CompileCommand struct {
	Directory string   // the working directory where the CompileCommand was executed
	Filename  string   // the filename associated with the CompileCommand
	Args      []string // Args[0] is the compiler executable
}

func newCompileCommand(c C.CXCompileCommand) CompileCommand {
	var r CompileCommand

	n := int(C.clang_CompileCommand_getNumArgs(c))

	r.Directory = cx2GoString(C.clang_CompileCommand_getDirectory(c))
	r.Filename = cx2GoString(C.clang_CompileCommand_getFilename(c))
	r.Args = make([]string, n)

	for i := range r.Args {
		r.Args[i] = cx2GoString(C.clang_CompileCommand_getArg(c, C.uint(i)))
	}
	return r
}
