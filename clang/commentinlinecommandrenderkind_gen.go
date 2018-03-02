package clang

// #include "./clang-c/Documentation.h"
// #include "go-clang.h"
import "C"

// The most appropriate rendering mode for an inline command, chosen on command semantics in Doxygen.
type CommentInlineCommandRenderKind uint32

const (
	// Command argument should be rendered in a normal font.
	CommentInlineCommandRenderKind_Normal CommentInlineCommandRenderKind = C.CXCommentInlineCommandRenderKind_Normal

	// Command argument should be rendered in a bold font.
	CommentInlineCommandRenderKind_Bold CommentInlineCommandRenderKind = C.CXCommentInlineCommandRenderKind_Bold

	// Command argument should be rendered in a monospaced font.
	CommentInlineCommandRenderKind_Monospaced CommentInlineCommandRenderKind = C.CXCommentInlineCommandRenderKind_Monospaced

	// Command argument should be rendered emphasized (typically italic font).
	CommentInlineCommandRenderKind_Emphasized CommentInlineCommandRenderKind = C.CXCommentInlineCommandRenderKind_Emphasized
)
