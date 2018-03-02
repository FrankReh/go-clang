package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

type IdxEntityLanguage uint32

const (
	IdxEntityLang_None  IdxEntityLanguage = C.CXIdxEntityLang_None
	IdxEntityLang_C     IdxEntityLanguage = C.CXIdxEntityLang_C
	IdxEntityLang_ObjC  IdxEntityLanguage = C.CXIdxEntityLang_ObjC
	IdxEntityLang_CXX   IdxEntityLanguage = C.CXIdxEntityLang_CXX
	IdxEntityLang_Swift IdxEntityLanguage = C.CXIdxEntityLang_Swift
)
