package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import "fmt"

type EvalResultKind uint32

const (
	Eval_Int            EvalResultKind = C.CXEval_Int
	Eval_Float          EvalResultKind = C.CXEval_Float
	Eval_ObjCStrLiteral EvalResultKind = C.CXEval_ObjCStrLiteral
	Eval_StrLiteral     EvalResultKind = C.CXEval_StrLiteral
	Eval_CFStr          EvalResultKind = C.CXEval_CFStr
	Eval_Other          EvalResultKind = C.CXEval_Other
	Eval_UnExposed      EvalResultKind = C.CXEval_UnExposed
)

func (erk EvalResultKind) String() string {
	switch erk {
	case Eval_Int:
		return "Eval_Int"
	case Eval_Float:
		return "Eval_Float"
	case Eval_ObjCStrLiteral:
		return "Eval_ObjCStrLiteral"
	case Eval_StrLiteral:
		return "Eval_StrLiteral"
	case Eval_CFStr:
		return "Eval_CFStr"
	case Eval_Other:
		return "Eval_Other"
	case Eval_UnExposed:
		return "Eval_UnExposed"
	}

	return fmt.Sprintf("EvalResultKind unknown %d", int(erk))
}
