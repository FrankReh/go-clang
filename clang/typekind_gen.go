package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import "github.com/frankreh/go-clang/clang/typekind"

func TypeKindSpelling(tk typekind.Kind) string {
	return cx2GoString(C.clang_getTypeKindSpelling(C.enum_CXTypeKind(tk)))
}
