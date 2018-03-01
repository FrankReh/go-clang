package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import "fmt"

type IndexOptFlags uint32

const (
	// Used to indicate that no special indexing options are needed.
	IndexOpt_None IndexOptFlags = C.CXIndexOpt_None

	// Used to indicate that IndexerCallbacks#indexEntityReference should be invoked for only one reference of an entity per source file that does not also include a declaration/definition of the entity.
	IndexOpt_SuppressRedundantRefs IndexOptFlags = C.CXIndexOpt_SuppressRedundantRefs

	// Function-local symbols should be indexed. If this is not set function-local symbols will be ignored.
	IndexOpt_IndexFunctionLocalSymbols IndexOptFlags = C.CXIndexOpt_IndexFunctionLocalSymbols

	// Implicit function/class template instantiations should be indexed. If this is not set, implicit instantiations will be ignored.
	IndexOpt_IndexImplicitTemplateInstantiations IndexOptFlags = C.CXIndexOpt_IndexImplicitTemplateInstantiations

	// Suppress all compiler warnings when parsing for indexing.
	IndexOpt_SuppressWarnings IndexOptFlags = C.CXIndexOpt_SuppressWarnings

	// Skip a function/method body that was already parsed during an indexing session associated with a CXIndexAction object. Bodies in system headers are always skipped.
	IndexOpt_SkipParsedBodiesInSession IndexOptFlags = C.CXIndexOpt_SkipParsedBodiesInSession
)

func (iof IndexOptFlags) String() string {
	switch iof {
	case IndexOpt_None:
		return "IndexOpt_None"
	case IndexOpt_SuppressRedundantRefs:
		return "IndexOpt_SuppressRedundantRefs"
	case IndexOpt_IndexFunctionLocalSymbols:
		return "IndexOpt_IndexFunctionLocalSymbols"
	case IndexOpt_IndexImplicitTemplateInstantiations:
		return "IndexOpt_IndexImplicitTemplateInstantiations"
	case IndexOpt_SuppressWarnings:
		return "IndexOpt_SuppressWarnings"
	case IndexOpt_SkipParsedBodiesInSession:
		return "IndexOpt_SkipParsedBodiesInSession"
	}

	return fmt.Sprintf("IndexOptFlags unknown %d", int(iof))
}
