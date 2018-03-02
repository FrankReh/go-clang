package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

// Describe the "language" of the entity referred to by a cursor.
type LanguageKind uint32

const (
	Language_Invalid   LanguageKind = C.CXLanguage_Invalid
	Language_C         LanguageKind = C.CXLanguage_C
	Language_ObjC      LanguageKind = C.CXLanguage_ObjC
	Language_CPlusPlus LanguageKind = C.CXLanguage_CPlusPlus
)
