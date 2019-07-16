#include "_cgo_export.h"
#include "go-clang.h"

unsigned go_clang_visit_children(CXCursor c, void *opaque) {
	return clang_visitChildren(c, (CXCursorVisitor)&GoClangCursorVisitor, opaque);
}
