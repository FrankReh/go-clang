package clang

// #include "go-clang.h"
import "C"

// A character string. The CXString type is used to return strings from the
// interface when the ownership of that string might differ from one call to
// the next. Use clang_getCString() to retrieve the string data and, once
// finished with the string data, call clang_disposeString() to free the
// string.
type cxstring struct {
	c C.CXString
}

func cx2GoString(c C.CXString) string {
	// Would be nice to create one C routine for getting the CString and
	// disposing of it right away.  Perhaps pass in a buffer of usually large
	// enough space.  When the buffer is not large enough, either call this old
	// way, or call with a larger buffer.
	cstr := C.clang_getCString(c)
	s := C.GoString(cstr)
	C.clang_disposeString(c)
	return s
}

// Retrieve the character data associated with the given string.
func (c cxstring) String() string {
	cstr := C.clang_getCString(c.c)
	return C.GoString(cstr)
}

// Free the given string.
func (c cxstring) Dispose() {
	C.clang_disposeString(c.c)
}
