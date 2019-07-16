/*
	type File struct {
		// Has unexported fields.
	}
		A particular source file that is part of a translation unit.

	func (f File) IsEqual(f2 File) bool
	func (f File) Name() string
	func (f File) Time() time.Time
	func (f File) TryGetRealPathName() string
	func (f File) UniqueID() (FileUniqueID, error)
*/
package clang

// #include "go-clang.h"
import "C"
import "time"

// A particular source file that is part of a translation unit.
type File struct{ c C.CXFile }

// Return the complete file and path name of the given file.
func (f File) Name() string { return cx2GoString(C.clang_getFileName(f.c)) }

// Returns the real path name of file.  An empty string may be returned.
// Use clang_getFileName() in that case.
func (f File) TryGetRealPathName() string { return cx2GoString(C.clang_File_tryGetRealPathName(f.c)) }

// Return the last modification time of the given file.
func (f File) Time() time.Time { return time.Unix(int64(C.clang_getFileTime(f.c)), 0) }

// Return the unique ID for the given file.
func (f File) UniqueID() (FileUniqueID, error) {
	var outID FileUniqueID
	var err error

	if o := int32(C.clang_getFileUniqueID(f.c, &outID.c)); o != 0 {
		err = UniqueIDErr
	}

	return outID, err
}

// The call to C.clang_getFileUniqueID returned nonzero, indicating an error occurred.
const UniqueIDErr = Error("UniqueID")

// Returns non-zero if the file1 and file2 point to the same file, or they are both NULL.
func (f File) IsEqual(f2 File) bool { return C.clang_File_isEqual(f.c, f2.c) != 0 }
