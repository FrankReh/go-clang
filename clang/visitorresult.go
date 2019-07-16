package clang

// #include "go-clang.h"
import "C"

type VisitorResult uint32

const (
	Visit_Break    VisitorResult = C.CXVisit_Break
	Visit_Continue VisitorResult = C.CXVisit_Continue
)
