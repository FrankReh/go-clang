package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import "fmt"

type IdxEntityLanguage uint32

const (
	IdxEntityLang_None  IdxEntityLanguage = C.CXIdxEntityLang_None
	IdxEntityLang_C     IdxEntityLanguage = C.CXIdxEntityLang_C
	IdxEntityLang_ObjC  IdxEntityLanguage = C.CXIdxEntityLang_ObjC
	IdxEntityLang_CXX   IdxEntityLanguage = C.CXIdxEntityLang_CXX
	IdxEntityLang_Swift IdxEntityLanguage = C.CXIdxEntityLang_Swift
)

func (iel IdxEntityLanguage) String() string {
	switch iel {
	case IdxEntityLang_None:
		return "IdxEntityLang_None"
	case IdxEntityLang_C:
		return "IdxEntityLang_C"
	case IdxEntityLang_ObjC:
		return "IdxEntityLang_ObjC"
	case IdxEntityLang_CXX:
		return "IdxEntityLang_CXX"
	case IdxEntityLang_Swift:
		return "IdxEntityLang_Swift"
	}

	return fmt.Sprintf("IdxEntityLanguage unkown %d", int(iel))
}
