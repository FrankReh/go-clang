package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import (
	"github.com/frankreh/go-clang/clang/printing"
)

/*
	Opaque pointer representing a policy that controls pretty printing for
	clang_getCursorPrettyPrinted.
*/
type PrintingPolicy struct {
	c C.CXPrintingPolicy
}

/**
 * Get a property value for the given printing policy.

	CINDEX_LINKAGE unsigned
	clang_PrintingPolicy_getProperty(CXPrintingPolicy Policy,
                                 enum CXPrintingPolicyProperty Property);
*/
func (p PrintingPolicy) Property(property printing.Property) uint {
	return uint(C.clang_PrintingPolicy_getProperty(p.c, uint32(property)))
}

/**
 * Set a property value for the given printing policy.

	CINDEX_LINKAGE void clang_PrintingPolicy_setProperty(CXPrintingPolicy Policy,
                                                     enum CXPrintingPolicyProperty Property,
                                                     unsigned Value);
*/
func (p PrintingPolicy) SetProperty(property printing.Property, value uint) {
	C.clang_PrintingPolicy_setProperty(p.c, uint32(property), C.uint(value))
}

/**
 * Retrieve the default policy for the cursor.
 *
 * The policy should be released after use with \c
 * clang_PrintingPolicy_dispose.
 */
func (c Cursor) PrintingPolicy() PrintingPolicy {
	return PrintingPolicy{C.clang_getCursorPrintingPolicy(c.c)}
}

/**
 * Release a printing policy.

	CINDEX_LINKAGE void clang_PrintingPolicy_dispose(CXPrintingPolicy Policy);
*/
func (p PrintingPolicy) Displose() {
	C.clang_PrintingPolicy_dispose(p.c)
}

/**
 * Pretty print declarations.
 *
 * Cursor The cursor representing a declaration.
 *
 * Policy The policy to control the entities being printed. If
 * NULL, a default policy is used.
 *
 * Returns The pretty printed declaration or the empty string for
 * other cursors.

	CINDEX_LINKAGE CXString clang_getCursorPrettyPrinted(CXCursor Cursor,
                                                     CXPrintingPolicy Policy);
*/
func (c Cursor) PrettyPrinted(policy PrintingPolicy) string {
	return cx2GoString(C.clang_getCursorPrettyPrinted(c.c, policy.c))
}
