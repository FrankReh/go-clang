package clang

// #include "go-clang.h"
import "C"

// Data for ppIncludedFile callback.
type IdxIncludedFileInfo struct {
	c C.CXIdxIncludedFileInfo
}

// Location of '#' in the \#include/\#import directive.
func (iifi IdxIncludedFileInfo) HashLoc() IdxLoc {
	return IdxLoc{iifi.c.hashLoc}
}

// Filename as written in the \#include/\#import directive.
func (iifi IdxIncludedFileInfo) Filename() string {
	return C.GoString(iifi.c.filename)
}

// The actual file that the \#include/\#import directive resolved to.
func (iifi IdxIncludedFileInfo) File() File {
	return File{iifi.c.file}
}

func (iifi IdxIncludedFileInfo) IsImport() bool {
	return iifi.c.isImport != 0
}

func (iifi IdxIncludedFileInfo) IsAngled() bool {
	return iifi.c.isAngled != 0
}

// Has the directive been automatically turned into a module import.
func (iifi IdxIncludedFileInfo) IsModuleImport() bool {
	return iifi.c.isModuleImport != 0
}
