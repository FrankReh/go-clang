package clang

// #include "go-clang.h"
import "C"

// A group of callbacks used by #clang_indexSourceFile and #clang_indexTranslationUnit.
type IndexerCallbacks struct{ c C.IndexerCallbacks }

//

// The client's data object that is associated with a CXFile.  Here some
// client's opaque go data could be pointed to.  Given the problems with go
// memory vs c, this would more than likely be an int that the go side would
// map to something more meaningful, like the visit callbacks.
//
// The client creates it with their two callback functions enteredMainFile
// and ppIncludedFile.
//
type IdxClientFile struct {
	c C.CXIdxClientFile
}

//

// Data for IndexerCallbacks#importedASTFile.
type IdxImportedASTFileInfo struct{ c C.CXIdxImportedASTFileInfo }

// Top level AST file containing the imported PCH, module or submodule.
func (g IdxImportedASTFileInfo) File() File { return File{g.c.file} }

// The imported module or NULL if the AST file is a PCH.
func (g IdxImportedASTFileInfo) Module() Module { return Module{g.c.module} }

// Location where the file is imported. Applicable only for modules.
func (g IdxImportedASTFileInfo) Loc() IdxLoc { return IdxLoc{g.c.loc} }

// Non-zero if an inclusion directive was automatically turned into a module import. Applicable only for modules.
func (g IdxImportedASTFileInfo) IsImplicit() bool { return g.c.isImplicit != 0 }
