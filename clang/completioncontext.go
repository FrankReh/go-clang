package clang

// #include "go-clang.h"
import "C"

/*
	Bits that represent the context under which completion is occurring.

	The enumerators in this enumeration may be bitwise-OR'd together if multiple
	contexts are occurring simultaneously.
*/
type CompletionContext uint32

const (
	// The context for completions is unexposed, as only Clang results should be included. (This is equivalent to having no context bits set.)
	CompletionContext_Unexposed CompletionContext = C.CXCompletionContext_Unexposed
	// Completions for any possible type should be included in the results.
	CompletionContext_AnyType CompletionContext = C.CXCompletionContext_AnyType
	// Completions for any possible value (variables, function calls, etc.) should be included in the results.
	CompletionContext_AnyValue CompletionContext = C.CXCompletionContext_AnyValue
	// Completions for values that resolve to an Objective-C object should be included in the results.
	CompletionContext_ObjCObjectValue CompletionContext = C.CXCompletionContext_ObjCObjectValue
	// Completions for values that resolve to an Objective-C selector should be included in the results.
	CompletionContext_ObjCSelectorValue CompletionContext = C.CXCompletionContext_ObjCSelectorValue
	// Completions for values that resolve to a C++ class type should be included in the results.
	CompletionContext_CXXClassTypeValue CompletionContext = C.CXCompletionContext_CXXClassTypeValue
	// Completions for fields of the member being accessed using the dot operator should be included in the results.
	CompletionContext_DotMemberAccess CompletionContext = C.CXCompletionContext_DotMemberAccess
	// Completions for fields of the member being accessed using the arrow operator should be included in the results.
	CompletionContext_ArrowMemberAccess CompletionContext = C.CXCompletionContext_ArrowMemberAccess
	// Completions for properties of the Objective-C object being accessed using the dot operator should be included in the results.
	CompletionContext_ObjCPropertyAccess CompletionContext = C.CXCompletionContext_ObjCPropertyAccess
	// Completions for enum tags should be included in the results.
	CompletionContext_EnumTag CompletionContext = C.CXCompletionContext_EnumTag
	// Completions for union tags should be included in the results.
	CompletionContext_UnionTag CompletionContext = C.CXCompletionContext_UnionTag
	// Completions for struct tags should be included in the results.
	CompletionContext_StructTag CompletionContext = C.CXCompletionContext_StructTag
	// Completions for C++ class names should be included in the results.
	CompletionContext_ClassTag CompletionContext = C.CXCompletionContext_ClassTag
	// Completions for C++ namespaces and namespace aliases should be included in the results.
	CompletionContext_Namespace CompletionContext = C.CXCompletionContext_Namespace
	// Completions for C++ nested name specifiers should be included in the results.
	CompletionContext_NestedNameSpecifier CompletionContext = C.CXCompletionContext_NestedNameSpecifier
	// Completions for Objective-C interfaces (classes) should be included in the results.
	CompletionContext_ObjCInterface CompletionContext = C.CXCompletionContext_ObjCInterface
	// Completions for Objective-C protocols should be included in the results.
	CompletionContext_ObjCProtocol CompletionContext = C.CXCompletionContext_ObjCProtocol
	// Completions for Objective-C categories should be included in the results.
	CompletionContext_ObjCCategory CompletionContext = C.CXCompletionContext_ObjCCategory
	// Completions for Objective-C instance messages should be included in the results.
	CompletionContext_ObjCInstanceMessage CompletionContext = C.CXCompletionContext_ObjCInstanceMessage
	// Completions for Objective-C class messages should be included in the results.
	CompletionContext_ObjCClassMessage CompletionContext = C.CXCompletionContext_ObjCClassMessage
	// Completions for Objective-C selector names should be included in the results.
	CompletionContext_ObjCSelectorName CompletionContext = C.CXCompletionContext_ObjCSelectorName
	// Completions for preprocessor macro names should be included in the results.
	CompletionContext_MacroName CompletionContext = C.CXCompletionContext_MacroName
	// Natural language completions should be included in the results.
	CompletionContext_NaturalLanguage CompletionContext = C.CXCompletionContext_NaturalLanguage
	// #include file completions should be included in the results.
	CompletionContext_IncludedFile CompletionContext = C.CXCompletionContext_IncludedFile
	// The current context is unknown, so set all contexts.
	CompletionContext_Unknown CompletionContext = C.CXCompletionContext_Unknown
)
