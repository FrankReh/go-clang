package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

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
