package clang

// #include "go-clang.h"
import "C"

// Represents the storage classes as declared in the source. CX_SC_Invalid was
// added for the case that the passed cursor in not a declaration.
type StorageClass uint32

const (
	SC_Invalid              StorageClass = C.CX_SC_Invalid
	SC_None                 StorageClass = C.CX_SC_None
	SC_Extern               StorageClass = C.CX_SC_Extern
	SC_Static               StorageClass = C.CX_SC_Static
	SC_PrivateExtern        StorageClass = C.CX_SC_PrivateExtern
	SC_OpenCLWorkGroupLocal StorageClass = C.CX_SC_OpenCLWorkGroupLocal
	SC_Auto                 StorageClass = C.CX_SC_Auto
	SC_Register             StorageClass = C.CX_SC_Register
)
