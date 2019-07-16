package clang

// #include "go-clang.h"
import "C"
import "unsafe"

// A remapping of original source files and their translated files.
type Remapping struct {
	Original, Transformed string
}

/*
	Retrieve a remapping.

	Parameter path the path that contains metadata about remappings.

	Returns the requested remapping. This remapping must be freed
	via a call to clang_remap_dispose(). Can return NULL if an error occurred.
*/
func NewRemappings(path string) []Remapping {
	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	return copyAndDisposeRemappings(C.clang_getRemappings(c_path))
}

/*
	Retrieve a remapping.

	Parameter filePaths pointer to an array of file paths containing remapping info.

	Parameter numFiles number of file paths.

	Returns the requested remapping. This remapping must be freed
	via a call to clang_remap_dispose(). Can return NULL if an error occurred.
*/
func NewRemappingsFromFileList(filePaths []string) []Remapping {
	ca_filePaths := make([]*C.char, len(filePaths))
	var cp_filePaths **C.char
	if len(filePaths) > 0 {
		cp_filePaths = &ca_filePaths[0]
	}
	for i := range filePaths {
		ci_str := C.CString(filePaths[i])
		defer C.free(unsafe.Pointer(ci_str))
		ca_filePaths[i] = ci_str
	}

	return copyAndDisposeRemappings(C.clang_getRemappingsFromFileList(cp_filePaths, C.uint(len(filePaths))))
}

func copyAndDisposeRemappings(c C.CXRemapping) []Remapping {
	if c == nil {
		return nil
	}

	n := int(C.clang_remap_getNumFiles(c))
	r := make([]Remapping, n)

	for i := range r {
		var original cxstring
		var transformed cxstring

		C.clang_remap_getFilenames(c, C.uint(i), &original.c, &transformed.c)

		r[i].Original = original.String()
		r[i].Transformed = transformed.String()

		original.Dispose()
		transformed.Dispose()
	}

	C.clang_remap_dispose(c)

	return r
}
