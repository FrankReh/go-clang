package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

type IdxDeclInfoFlags uint32

const (
	IdxDeclFlag_Skipped IdxDeclInfoFlags = C.CXIdxDeclFlag_Skipped
)
