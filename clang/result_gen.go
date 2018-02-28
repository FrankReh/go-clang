package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import "fmt"

type Result uint32

const (
	// Function returned successfully.
	Result_Success Result = C.CXResult_Success

	// One of the parameters was invalid for the function.
	Result_Invalid Result = C.CXResult_Invalid

	// The function was terminated by a callback (e.g. it returned CXVisit_Break)
	Result_VisitBreak Result = C.CXResult_VisitBreak
)

func (r Result) String() string {
	switch r {
	case Result_Success:
		return "Result_Success"
	case Result_Invalid:
		return "Result_Invalid"
	case Result_VisitBreak:
		return "Result_VisitBreak"
	}

	return fmt.Sprintf("Result unknown %d", int(r))
}
