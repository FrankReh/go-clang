package clang

// #include "./clang-c/Documentation.h"
// #include "go-clang.h"
import "C"

// Describes parameter passing direction for Parameter or arg command.
type CommentParamPassDirection uint32

const (
	// The parameter is an input parameter.
	CommentParamPassDirection_In CommentParamPassDirection = C.CXCommentParamPassDirection_In

	// The parameter is an output parameter.
	CommentParamPassDirection_Out CommentParamPassDirection = C.CXCommentParamPassDirection_Out

	// The parameter is an input and output parameter.
	CommentParamPassDirection_InOut CommentParamPassDirection = C.CXCommentParamPassDirection_InOut
)
