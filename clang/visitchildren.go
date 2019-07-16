package clang

// #include "go-clang.h"
import "C"
import (
	"sync"
	"unsafe"
)

// CursorVisitor is the callback function type passed to Visit.
/**
 * Visitor invoked for each cursor found by a traversal.
 *
 * This visitor function will be invoked for each cursor found by
 * clang_visitCursorChildren(). Its first argument is the cursor being
 * visited, its second argument is the parent visitor for that cursor,
 * and its third argument is the client data provided to
 * clang_visitCursorChildren().
 *
 * The visitor should return one of the CXChildVisitResult values
 * to direct clang_visitCursorChildren().
 */
type CursorVisitor func(cursor, parent Cursor) ChildVisitResult

type funcRegistry struct {
	sync.RWMutex

	index uintptr // Gets passed to C as (void*) via unsafe.Pointer, so uintptr is most convenient.
	funcs map[uintptr]CursorVisitor
}

func (fm *funcRegistry) register(f CursorVisitor) uintptr {
	fm.Lock()
	defer fm.Unlock()

	fm.index++
	for fm.funcs[fm.index] != nil {
		fm.index++
	}

	fm.funcs[fm.index] = f

	return fm.index
}

func (fm *funcRegistry) lookup(index uintptr) CursorVisitor {
	fm.RLock()
	defer fm.RUnlock()

	return fm.funcs[index]
}

func (fm *funcRegistry) unregister(index uintptr) {
	fm.Lock()

	delete(fm.funcs, index)

	fm.Unlock()
}

var visitors = funcRegistry{
	funcs: make(map[uintptr]CursorVisitor),
}

// Visit invokes the visitor callback on a cursor's children.
// Probably a misnomer. Should have been named VisitChildren.
/**
 * This function visits all the direct children of the given cursor,
 * invoking the given visitor function with the cursors of each
 * visited child. The traversal may be recursive, if the visitor returns
 * CXChildVisit_Recurse. The traversal may also be ended prematurely, if
 * the visitor returns CXChildVisit_Break.
 *
 * param parent the cursor whose child may be visited. All kinds of
 * cursors can be visited, including invalid cursors (which, by
 * definition, have no children).
 *
 * param visitor the visitor function that will be invoked for each
 * child of parent.
 *
 * param client_data pointer data supplied by the client, which will
 * be passed to the visitor each time it is invoked.
 *
 * returns a non-zero value if the traversal was terminated
 * prematurely by the visitor returning CXChildVisit_Break.
 */
func (c Cursor) Visit(visitor CursorVisitor) bool {
	index := visitors.register(visitor)
	defer visitors.unregister(index)

	o := C.go_clang_visit_children(c.c, unsafe.Pointer(index))

	return o == 0
}

// GoClangCursorVisitor calls the cursor visitor
//export GoClangCursorVisitor
func GoClangCursorVisitor(cursor C.CXCursor, parent C.CXCursor, opaque unsafe.Pointer) ChildVisitResult {
	index := uintptr(opaque)
	fn := visitors.lookup(index)

	if fn == nil {
		// TODO consider calling panic or setting an error.
		return ChildVisit_Break
	}

	return fn(Cursor{cursor}, Cursor{parent})
}

// Describes how the traversal of the children of a particular cursor should
// proceed after visiting a particular child cursor.
//
// A value of this enumeration type should be returned by each CXCursorVisitor
// to indicate how clang_visitChildren() proceed.
type ChildVisitResult uint32

const (
	// Terminates the cursor traversal.
	ChildVisit_Break ChildVisitResult = C.CXChildVisit_Break

	// Continues the cursor traversal with the next sibling of the cursor just visited, without visiting its children.
	ChildVisit_Continue ChildVisitResult = C.CXChildVisit_Continue

	// Recursively traverse the children of this cursor, using the same visitor and client data.
	ChildVisit_Recurse ChildVisitResult = C.CXChildVisit_Recurse
)
