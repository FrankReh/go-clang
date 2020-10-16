package clang

// #include "go-clang.h"
import "C"

/*
	Flags that control the creation of translation units.

	The enumerators in this enumeration type are meant to be bitwise
	ORed together to specify which options should be used when
	constructing the translation unit.
*/
type TranslationUnit_Flags uint32

const (
	/*
		Used to indicate that no special translation-unit options are needed.
	*/
	TranslationUnit_None TranslationUnit_Flags = C.CXTranslationUnit_None

	/*
		Used to indicate that the parser should construct a "detailed"
		preprocessing record, including all macro definitions and instantiations.

		Constructing a detailed preprocessing record requires more memory
		and time to parse, since the information contained in the record
		is usually not retained. However, it can be useful for
		applications that require more detailed information about the
		behavior of the preprocessor.
	*/
	TranslationUnit_DetailedPreprocessingRecord TranslationUnit_Flags = C.CXTranslationUnit_DetailedPreprocessingRecord

	/*
		Used to indicate that the translation unit is incomplete.

		When a translation unit is considered "incomplete", semantic
		analysis that is typically performed at the end of the
		translation unit will be suppressed. For example, this suppresses
		the completion of tentative declarations in C and of
		instantiation of implicitly-instantiation function templates in
		C++. This option is typically used when parsing a header with the
		intent of producing a precompiled header.
	*/
	TranslationUnit_Incomplete TranslationUnit_Flags = C.CXTranslationUnit_Incomplete

	/*
		Used to indicate that the translation unit should be built with an
		implicit precompiled header for the preamble.

		An implicit precompiled header is used as an optimization when a
		particular translation unit is likely to be reparsed many times
		when the sources aren't changing that often. In this case, an
		implicit precompiled header will be built containing all of the
		initial includes at the top of the main file (what we refer to as
		the "preamble" of the file). In subsequent parses, if the
		preamble or the files in it have not changed,
		clang_reparseTranslationUnit() will re-use the implicit
		precompiled header to improve parsing performance.
	*/
	TranslationUnit_PrecompiledPreamble TranslationUnit_Flags = C.CXTranslationUnit_PrecompiledPreamble

	/*
		Used to indicate that the translation unit should cache some
		code-completion results with each reparse of the source file.

		Caching of code-completion results is a performance optimization that
		introduces some overhead to reparsing but improves the performance of
		code-completion operations.
	*/
	TranslationUnit_CacheCompletionResults TranslationUnit_Flags = C.CXTranslationUnit_CacheCompletionResults

	/*
		Used to indicate that the translation unit will be serialized with
		clang_saveTranslationUnit.

		This option is typically used when parsing a header with the intent of
		producing a precompiled header.
	*/
	TranslationUnit_ForSerialization TranslationUnit_Flags = C.CXTranslationUnit_ForSerialization

	/*
		DEPRECATED: Enabled chained precompiled preambles in C++.

		Note: this is a *temporary* option that is available only while
		we are testing C++ precompiled preamble support. It is deprecated.
	*/
	TranslationUnit_CXXChainedPCH TranslationUnit_Flags = C.CXTranslationUnit_CXXChainedPCH

	/*
		Used to indicate that function/method bodies should be skipped while
		parsing.

		This option can be used to search for declarations/definitions while
		ignoring the usages.
	*/
	TranslationUnit_SkipFunctionBodies TranslationUnit_Flags = C.CXTranslationUnit_SkipFunctionBodies

	/*
		Used to indicate that brief documentation comments should be included into the set of code completions returned from this translation unit.
	*/
	TranslationUnit_IncludeBriefCommentsInCodeCompletion TranslationUnit_Flags = C.CXTranslationUnit_IncludeBriefCommentsInCodeCompletion

	/*
		Used to indicate that the precompiled preamble should be created on the first parse. Otherwise it will be created on the first reparse. This trades runtime on the first parse (serializing the preamble takes time) for reduced runtime on the second parse (can now reuse the preamble).
	*/
	TranslationUnit_CreatePreambleOnFirstParse TranslationUnit_Flags = C.CXTranslationUnit_CreatePreambleOnFirstParse

	/*
		Do not stop processing when fatal errors are encountered.

		When fatal errors are encountered while parsing a translation unit,
		semantic analysis is typically stopped early when compiling code. A common
		source for fatal errors are unresolvable include files. For the
		purposes of an IDE, this is undesirable behavior and as much information
		as possible should be reported. Use this flag to enable this behavior.
	*/
	TranslationUnit_KeepGoing TranslationUnit_Flags = C.CXTranslationUnit_KeepGoing

	/*
		Sets the preprocessor in a mode for parsing a single file only.
	*/
	TranslationUnit_SingleFileParse TranslationUnit_Flags = C.CXTranslationUnit_SingleFileParse

	/*
		Used in combination with CXTranslationUnit_SkipFunctionBodies to
		constrain the skipping of function bodies to the preamble.

		The function bodies of the main file are not skipped.
	*/

	TranslationUnit_LimitSkipFunctionBodiesToPreamble TranslationUnit_Flags = C.CXTranslationUnit_LimitSkipFunctionBodiesToPreamble

	/*
		Used to indicate that attributed types should be included in CXType.
	*/
	TranslationUnit_IncludeAttributedTypes TranslationUnit_Flags = C.CXTranslationUnit_IncludeAttributedTypes

	/*
		Used to indicate that implicit attributes should be visited.
	*/
	TranslationUnit_VisitImplicitAttributes TranslationUnit_Flags = C.CXTranslationUnit_VisitImplicitAttributes

	/*
		Used to indicate that non-errors from included files should be ignored.
	*/
	TranslationUnit_IgnoreNonErrorsFromIncludedFiles TranslationUnit_Flags = C.CXTranslationUnit_IgnoreNonErrorsFromIncludedFiles

	/*
		Tells the preprocessor not to skip excluded conditional blocks.
	*/
	TranslationUnit_RetainExcludedConditionalBlocks TranslationUnit_Flags = C.CXTranslationUnit_RetainExcludedConditionalBlocks
)
