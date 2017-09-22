package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import (
	"reflect"
	"unsafe"
)

type ResourceUsageEntry struct {
	kind   TUResourceUsageKind
	amount uint64
}

func (rue ResourceUsageEntry) Kind() TUResourceUsageKind {
	return rue.kind
}

func (rue ResourceUsageEntry) Amount() uint64 {
	return rue.amount
}

// Return the memory usage of a translation unit. This object should be released with clang_disposeCXTUResourceUsage().
func (tu TranslationUnit) ResourceUsage() []ResourceUsageEntry {
	ru := C.clang_getCXTUResourceUsage(tu.c)
	defer C.clang_disposeCXTUResourceUsage(ru)

	// Build a temporary slice with a shared backing store long enough to build
	// a normal one from its entries.
	var s []C.CXTUResourceUsageEntry
	gos_s := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	gos_s.Cap = int(ru.numEntries)
	gos_s.Len = int(ru.numEntries)
	gos_s.Data = uintptr(unsafe.Pointer(ru.entries))

	r := make([]ResourceUsageEntry, len(s))
	for i := range r {
		r[i].kind = TUResourceUsageKind(s[i].kind)
		r[i].amount = uint64(s[i].amount)
	}

	return r
}
