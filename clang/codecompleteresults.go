package clang

// #include "go-clang.h"
import "C"
import (
	"reflect"
	"unsafe"

	"github.com/frankreh/go-clang/clang/cursorkind"
)

// Contains the results of code-completion.
//
// This data structure contains the results of code completion, as produced by
// clang_codeCompleteAt(). Its contents must be freed by clang_disposeCodeCompleteResults.
type CodeCompleteResults struct {
	c *C.CXCodeCompleteResults
}

// Free the given set of code-completion results.
func (ccr *CodeCompleteResults) Dispose() {
	C.clang_disposeCodeCompleteResults(ccr.c)
}

/*
	Retrieve the number of fix-its for the given completion index.

	Calling this makes sense only if CXCodeComplete_IncludeCompletionsWithFixIts
	option was set.

	results The structure keeping all completion results

	completion_index The index of the completion

	Return The number of fix-its which must be applied before the completion at
	completion_index can be applied
*/
func (ccr *CodeCompleteResults) NumFixItsFor(completion_index uint) uint {
	return uint(C.clang_getCompletionNumFixIts(ccr.c, C.uint(completion_index)))
}

/*
	Fix-its that *must* be applied before inserting the text for the
	corresponding completion.

	By default, clang_codeCompleteAt() only returns completions with empty
	fix-its. Extra completions with non-empty fix-its should be explicitly
	requested by setting CXCodeComplete_IncludeCompletionsWithFixIts.

	For the clients to be able to compute position of the cursor after applying
	fix-its, the following conditions are guaranteed to hold for
	replacement_range of the stored fix-its:
	- Ranges in the fix-its are guaranteed to never contain the completion
	  point (or identifier under completion point, if any) inside them, except
	  at the start or at the end of the range.
	- If a fix-it range starts or ends with completion point (or starts or
	  ends after the identifier under completion point), it will contain at
	  least one character. It allows to unambiguously recompute completion
	  point after applying the fix-it.

	The intuition is that provided fix-its change code around the identifier we
	complete, but are not allowed to touch the identifier itself or the
	completion point. One example of completions with corrections are the ones
	replacing '.' with '->' and vice versa:

	std::unique_ptr<std::vector<int>> vec_ptr;
	In 'vec_ptr.^', one of the completions is 'push_back', it requires
	replacing '.' with '->'.
	In 'vec_ptr->^', one of the completions is 'release', it requires
	replacing '->' with '.'.

	results The structure keeping all completion results

	completion_index The index of the completion

	fixit_index The index of the fix-it for the completion at
	completion_index

	replacement_range The fix-it range that must be replaced before the
	completion at completion_index can be applied

	Returns The fix-it string that must replace the code at replacement_range
	before the completion at completion_index can be applied
*/
func (ccr *CodeCompleteResults) FixIt(completion_index, fixit_index uint, replacement_range SourceRange) string {
	return cx2GoString(C.clang_getCompletionFixIt(ccr.c, C.uint(completion_index), C.uint(fixit_index), &replacement_range.c))
}

// Determine the number of diagnostics produced prior to the location where code completion was performed.
func (ccr *CodeCompleteResults) NumDiagnostics() uint32 {
	return uint32(C.clang_codeCompleteGetNumDiagnostics(ccr.c))
}

/*
	Retrieve a diagnostic associated with the given code completion.

	Parameter Results the code completion results to query.
	Parameter Index the zero-based diagnostic number to retrieve.

	Returns the requested diagnostic. This diagnostic must be freed
	via a call to clang_disposeDiagnostic().
*/
func (ccr *CodeCompleteResults) Diagnostic(index uint32) Diagnostic {
	return Diagnostic{C.clang_codeCompleteGetDiagnostic(ccr.c, C.uint(index))}
}

/*
	Determines what completions are appropriate for the context
	the given code completion.

	Parameter Results the code completion results to query

	Returns the kinds of completions that are appropriate for use
	along with the given code completion results.
*/
func (ccr *CodeCompleteResults) Contexts() uint64 {
	return uint64(C.clang_codeCompleteGetContexts(ccr.c))
}

/*
	Returns the cursor kind for the container for the current code
	completion context. The container is only guaranteed to be set for
	contexts where a container exists (i.e. member accesses or Objective-C
	message sends); if there is not a container, this function will return
	CXCursor_InvalidCode.

	Parameter Results the code completion results to query

	Parameter IsIncomplete on return, this value will be false if Clang has complete
	information about the container. If Clang does not have complete
	information, this value will be true.

	Returns the container kind, or CXCursor_InvalidCode if there is not a
	container
*/
func (ccr *CodeCompleteResults) ContainerKind() (uint32, cursorkind.Kind) {
	var isIncomplete C.uint

	o := cursorkind.MustValidate(int(C.clang_codeCompleteGetContainerKind(ccr.c, &isIncomplete)))

	return uint32(isIncomplete), o // TBD swap and make this an error
}

/*
	Returns the USR for the container for the current code completion
	context. If there is not a container for the current context, this
	function will return the empty string.

	Parameter Results the code completion results to query

	Returns the USR for the container
*/
func (ccr *CodeCompleteResults) ContainerUSR() string {
	return cx2GoString(C.clang_codeCompleteGetContainerUSR(ccr.c))
}

/*
	Returns the currently-entered selector for an Objective-C message
	send, formatted like "initWithFoo:bar:". Only guaranteed to return a
	non-empty string for CXCompletionContext_ObjCInstanceMessage and
	CXCompletionContext_ObjCClassMessage.

	Parameter Results the code completion results to query

	Returns the selector (or partial selector) that has been entered thus far
	for an Objective-C message send.
*/
func (ccr *CodeCompleteResults) Selector() string {
	return cx2GoString(C.clang_codeCompleteGetObjCSelector(ccr.c))
}

// The code-completion results.
func (ccr CodeCompleteResults) Results() []CompletionResult {
	var s []CompletionResult
	gos_s := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	gos_s.Cap = int(ccr.c.NumResults)
	gos_s.Len = int(ccr.c.NumResults)
	gos_s.Data = uintptr(unsafe.Pointer(ccr.c.Results))

	// Create a slice with a backing store not shared with C to return to user,
	// allowing the ccr to be disposed and the slice returned here to live on.
	r := make([]CompletionResult, len(s))
	copy(r, s)

	return r
}

// The number of code-completion results stored in the Results array.
func (ccr CodeCompleteResults) NumResults() uint32 {
	return uint32(ccr.c.NumResults)
}

func (ccr *CodeCompleteResults) Diagnostics() []Diagnostic {
	s := make([]Diagnostic, ccr.NumDiagnostics())

	for i := range s {
		s[i] = ccr.Diagnostic(uint32(i))
	}

	return s
}
