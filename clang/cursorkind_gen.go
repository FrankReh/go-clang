package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

// Describes the kind of entity that a cursor refers to.
type CursorKind uint32

const (
	/*
		A declaration whose specific kind is not exposed via this
		interface.

		Unexposed declarations have the same operations as any other kind
		of declaration; one can extract their location information,
		spelling, find their definitions, etc. However, the specific kind
		of the declaration is not reported.
	*/
	Cursor_UnexposedDecl CursorKind = C.CXCursor_UnexposedDecl
	// A C or C++ struct.
	Cursor_StructDecl CursorKind = C.CXCursor_StructDecl
	// A C or C++ union.
	Cursor_UnionDecl CursorKind = C.CXCursor_UnionDecl
	// A C++ class.
	Cursor_ClassDecl CursorKind = C.CXCursor_ClassDecl
	// An enumeration.
	Cursor_EnumDecl CursorKind = C.CXCursor_EnumDecl
	// A field (in C) or non-static data member (in C++) in a struct, union, or C++ class.
	Cursor_FieldDecl CursorKind = C.CXCursor_FieldDecl
	// An enumerator constant.
	Cursor_EnumConstantDecl CursorKind = C.CXCursor_EnumConstantDecl
	// A function.
	Cursor_FunctionDecl CursorKind = C.CXCursor_FunctionDecl
	// A variable.
	Cursor_VarDecl CursorKind = C.CXCursor_VarDecl
	// A function or method parameter.
	Cursor_ParmDecl CursorKind = C.CXCursor_ParmDecl
	// An Objective-C \@interface.
	Cursor_ObjCInterfaceDecl CursorKind = C.CXCursor_ObjCInterfaceDecl
	// An Objective-C \@interface for a category.
	Cursor_ObjCCategoryDecl CursorKind = C.CXCursor_ObjCCategoryDecl
	// An Objective-C \@protocol declaration.
	Cursor_ObjCProtocolDecl CursorKind = C.CXCursor_ObjCProtocolDecl
	// An Objective-C \@property declaration.
	Cursor_ObjCPropertyDecl CursorKind = C.CXCursor_ObjCPropertyDecl
	// An Objective-C instance variable.
	Cursor_ObjCIvarDecl CursorKind = C.CXCursor_ObjCIvarDecl
	// An Objective-C instance method.
	Cursor_ObjCInstanceMethodDecl CursorKind = C.CXCursor_ObjCInstanceMethodDecl
	// An Objective-C class method.
	Cursor_ObjCClassMethodDecl CursorKind = C.CXCursor_ObjCClassMethodDecl
	// An Objective-C \@implementation.
	Cursor_ObjCImplementationDecl CursorKind = C.CXCursor_ObjCImplementationDecl
	// An Objective-C \@implementation for a category.
	Cursor_ObjCCategoryImplDecl CursorKind = C.CXCursor_ObjCCategoryImplDecl
	// A typedef.
	Cursor_TypedefDecl CursorKind = C.CXCursor_TypedefDecl
	// A C++ class method.
	Cursor_CXXMethod CursorKind = C.CXCursor_CXXMethod
	// A C++ namespace.
	Cursor_Namespace CursorKind = C.CXCursor_Namespace
	// A linkage specification, e.g. 'extern "C"'.
	Cursor_LinkageSpec CursorKind = C.CXCursor_LinkageSpec
	// A C++ constructor.
	Cursor_Constructor CursorKind = C.CXCursor_Constructor
	// A C++ destructor.
	Cursor_Destructor CursorKind = C.CXCursor_Destructor
	// A C++ conversion function.
	Cursor_ConversionFunction CursorKind = C.CXCursor_ConversionFunction
	// A C++ template type parameter.
	Cursor_TemplateTypeParameter CursorKind = C.CXCursor_TemplateTypeParameter
	// A C++ non-type template parameter.
	Cursor_NonTypeTemplateParameter CursorKind = C.CXCursor_NonTypeTemplateParameter
	// A C++ template template parameter.
	Cursor_TemplateTemplateParameter CursorKind = C.CXCursor_TemplateTemplateParameter
	// A C++ function template.
	Cursor_FunctionTemplate CursorKind = C.CXCursor_FunctionTemplate
	// A C++ class template.
	Cursor_ClassTemplate CursorKind = C.CXCursor_ClassTemplate
	// A C++ class template partial specialization.
	Cursor_ClassTemplatePartialSpecialization CursorKind = C.CXCursor_ClassTemplatePartialSpecialization
	// A C++ namespace alias declaration.
	Cursor_NamespaceAlias CursorKind = C.CXCursor_NamespaceAlias
	// A C++ using directive.
	Cursor_UsingDirective CursorKind = C.CXCursor_UsingDirective
	// A C++ using declaration.
	Cursor_UsingDeclaration CursorKind = C.CXCursor_UsingDeclaration
	// A C++ alias declaration
	Cursor_TypeAliasDecl CursorKind = C.CXCursor_TypeAliasDecl
	// An Objective-C \@synthesize definition.
	Cursor_ObjCSynthesizeDecl CursorKind = C.CXCursor_ObjCSynthesizeDecl
	// An Objective-C \@dynamic definition.
	Cursor_ObjCDynamicDecl CursorKind = C.CXCursor_ObjCDynamicDecl
	// An access specifier.
	Cursor_CXXAccessSpecifier CursorKind = C.CXCursor_CXXAccessSpecifier
	// An access specifier.
	Cursor_ObjCSuperClassRef CursorKind = C.CXCursor_ObjCSuperClassRef
	// An access specifier.
	Cursor_ObjCProtocolRef CursorKind = C.CXCursor_ObjCProtocolRef
	// An access specifier.
	Cursor_ObjCClassRef CursorKind = C.CXCursor_ObjCClassRef
	/*
		A reference to a type declaration.

		A type reference occurs anywhere where a type is named but not
		declared. For example, given:

		\code
		typedef unsigned size_type;
		size_type size;
		\endcode

		The typedef is a declaration of size_type (CXCursor_TypedefDecl),
		while the type of the variable "size" is referenced. The cursor
		referenced by the type of size is the typedef for size_type.
	*/
	Cursor_TypeRef CursorKind = C.CXCursor_TypeRef
	/*
		A reference to a type declaration.

		A type reference occurs anywhere where a type is named but not
		declared. For example, given:

		\code
		typedef unsigned size_type;
		size_type size;
		\endcode

		The typedef is a declaration of size_type (CXCursor_TypedefDecl),
		while the type of the variable "size" is referenced. The cursor
		referenced by the type of size is the typedef for size_type.
	*/
	Cursor_CXXBaseSpecifier CursorKind = C.CXCursor_CXXBaseSpecifier
	// A reference to a class template, function template, template template parameter, or class template partial specialization.
	Cursor_TemplateRef CursorKind = C.CXCursor_TemplateRef
	// A reference to a namespace or namespace alias.
	Cursor_NamespaceRef CursorKind = C.CXCursor_NamespaceRef
	// A reference to a member of a struct, union, or class that occurs in some non-expression context, e.g., a designated initializer.
	Cursor_MemberRef CursorKind = C.CXCursor_MemberRef
	/*
		A reference to a labeled statement.

		This cursor kind is used to describe the jump to "start_over" in the
		goto statement in the following example:

		\code
		start_over:
		++counter;

		goto start_over;
		\endcode

		A label reference cursor refers to a label statement.
	*/
	Cursor_LabelRef CursorKind = C.CXCursor_LabelRef
	/*
		A reference to a set of overloaded functions or function templates
		that has not yet been resolved to a specific function or function template.

		An overloaded declaration reference cursor occurs in C++ templates where
		a dependent name refers to a function. For example:

		\code
		template<typename T> void swap(T&, T&);

		struct X { ... };
		void swap(X&, X&);

		template<typename T>
		void reverse(T* first, T* last) {
		while (first < last - 1) {
		swap(*first, *--last);
		++first;
		}
		}

		struct Y { };
		void swap(Y&, Y&);
		\endcode

		Here, the identifier "swap" is associated with an overloaded declaration
		reference. In the template definition, "swap" refers to either of the two
		"swap" functions declared above, so both results will be available. At
		instantiation time, "swap" may also refer to other functions found via
		argument-dependent lookup (e.g., the "swap" function at the end of the
		example).

		The functions clang_getNumOverloadedDecls() and
		clang_getOverloadedDecl() can be used to retrieve the definitions
		referenced by this cursor.
	*/
	Cursor_OverloadedDeclRef CursorKind = C.CXCursor_OverloadedDeclRef
	// A reference to a variable that occurs in some non-expression context, e.g., a C++ lambda capture list.
	Cursor_VariableRef CursorKind = C.CXCursor_VariableRef
	// A reference to a variable that occurs in some non-expression context, e.g., a C++ lambda capture list.
	Cursor_InvalidFile CursorKind = C.CXCursor_InvalidFile
	// A reference to a variable that occurs in some non-expression context, e.g., a C++ lambda capture list.
	Cursor_NoDeclFound CursorKind = C.CXCursor_NoDeclFound
	// A reference to a variable that occurs in some non-expression context, e.g., a C++ lambda capture list.
	Cursor_NotImplemented CursorKind = C.CXCursor_NotImplemented
	// A reference to a variable that occurs in some non-expression context, e.g., a C++ lambda capture list.
	Cursor_InvalidCode CursorKind = C.CXCursor_InvalidCode
	/*
		An expression whose specific kind is not exposed via this
		interface.

		Unexposed expressions have the same operations as any other kind
		of expression; one can extract their location information,
		spelling, children, etc. However, the specific kind of the
		expression is not reported.
	*/
	Cursor_UnexposedExpr CursorKind = C.CXCursor_UnexposedExpr
	// An expression that refers to some value declaration, such as a function, variable, or enumerator.
	Cursor_DeclRefExpr CursorKind = C.CXCursor_DeclRefExpr
	// An expression that refers to a member of a struct, union, class, Objective-C class, etc.
	Cursor_MemberRefExpr CursorKind = C.CXCursor_MemberRefExpr
	// An expression that calls a function.
	Cursor_CallExpr CursorKind = C.CXCursor_CallExpr
	// An expression that sends a message to an Objective-C  object or class.
	Cursor_ObjCMessageExpr CursorKind = C.CXCursor_ObjCMessageExpr
	// An expression that represents a block literal.
	Cursor_BlockExpr CursorKind = C.CXCursor_BlockExpr
	// An integer literal.
	Cursor_IntegerLiteral CursorKind = C.CXCursor_IntegerLiteral
	// A floating point number literal.
	Cursor_FloatingLiteral CursorKind = C.CXCursor_FloatingLiteral
	// An imaginary number literal.
	Cursor_ImaginaryLiteral CursorKind = C.CXCursor_ImaginaryLiteral
	// A string literal.
	Cursor_StringLiteral CursorKind = C.CXCursor_StringLiteral
	// A character literal.
	Cursor_CharacterLiteral CursorKind = C.CXCursor_CharacterLiteral
	/*
		A parenthesized expression, e.g. "(1)".

		This AST node is only formed if full location information is requested.
	*/
	Cursor_ParenExpr CursorKind = C.CXCursor_ParenExpr
	// This represents the unary-expression's (except sizeof and alignof).
	Cursor_UnaryOperator CursorKind = C.CXCursor_UnaryOperator
	// [C99 6.5.2.1] Array Subscripting.
	Cursor_ArraySubscriptExpr CursorKind = C.CXCursor_ArraySubscriptExpr
	// A builtin binary operation expression such as "x + y" or "x <= y".
	Cursor_BinaryOperator CursorKind = C.CXCursor_BinaryOperator
	// Compound assignment such as "+=".
	Cursor_CompoundAssignOperator CursorKind = C.CXCursor_CompoundAssignOperator
	// The ?: ternary operator.
	Cursor_ConditionalOperator CursorKind = C.CXCursor_ConditionalOperator
	/*
		An explicit cast in C (C99 6.5.4) or a C-style cast in C++
		(C++ [expr.cast]), which uses the syntax (Type)expr.

		For example: (int)f.
	*/
	Cursor_CStyleCastExpr CursorKind = C.CXCursor_CStyleCastExpr
	// [C99 6.5.2.5]
	Cursor_CompoundLiteralExpr CursorKind = C.CXCursor_CompoundLiteralExpr
	// Describes an C or C++ initializer list.
	Cursor_InitListExpr CursorKind = C.CXCursor_InitListExpr
	// The GNU address of label extension, representing &&label.
	Cursor_AddrLabelExpr CursorKind = C.CXCursor_AddrLabelExpr
	// This is the GNU Statement Expression extension: ({int X=4; X;})
	Cursor_StmtExpr CursorKind = C.CXCursor_StmtExpr
	// Represents a C11 generic selection.
	Cursor_GenericSelectionExpr CursorKind = C.CXCursor_GenericSelectionExpr
	/*
		Implements the GNU __null extension, which is a name for a null
		pointer constant that has integral type (e.g., int or long) and is the same
		size and alignment as a pointer.

		The __null extension is typically only used by system headers, which define
		NULL as __null in C++ rather than using 0 (which is an integer that may not
		match the size of a pointer).
	*/
	Cursor_GNUNullExpr CursorKind = C.CXCursor_GNUNullExpr
	// C++'s static_cast<> expression.
	Cursor_CXXStaticCastExpr CursorKind = C.CXCursor_CXXStaticCastExpr
	// C++'s dynamic_cast<> expression.
	Cursor_CXXDynamicCastExpr CursorKind = C.CXCursor_CXXDynamicCastExpr
	// C++'s reinterpret_cast<> expression.
	Cursor_CXXReinterpretCastExpr CursorKind = C.CXCursor_CXXReinterpretCastExpr
	// C++'s const_cast<> expression.
	Cursor_CXXConstCastExpr CursorKind = C.CXCursor_CXXConstCastExpr
	/*
		Represents an explicit C++ type conversion that uses "functional"
		notion (C++ [expr.type.conv]).

		Example:
		\code
		x = int(0.5);
		\endcode
	*/
	Cursor_CXXFunctionalCastExpr CursorKind = C.CXCursor_CXXFunctionalCastExpr
	// A C++ typeid expression (C++ [expr.typeid]).
	Cursor_CXXTypeidExpr CursorKind = C.CXCursor_CXXTypeidExpr
	// [C++ 2.13.5] C++ Boolean Literal.
	Cursor_CXXBoolLiteralExpr CursorKind = C.CXCursor_CXXBoolLiteralExpr
	// [C++0x 2.14.7] C++ Pointer Literal.
	Cursor_CXXNullPtrLiteralExpr CursorKind = C.CXCursor_CXXNullPtrLiteralExpr
	// Represents the "this" expression in C++
	Cursor_CXXThisExpr CursorKind = C.CXCursor_CXXThisExpr
	/*
		[C++ 15] C++ Throw Expression.

		This handles 'throw' and 'throw' assignment-expression. When
		assignment-expression isn't present, Op will be null.
	*/
	Cursor_CXXThrowExpr CursorKind = C.CXCursor_CXXThrowExpr
	// A new expression for memory allocation and constructor calls, e.g: "new CXXNewExpr(foo)".
	Cursor_CXXNewExpr CursorKind = C.CXCursor_CXXNewExpr
	// A delete expression for memory deallocation and destructor calls, e.g. "delete[] pArray".
	Cursor_CXXDeleteExpr CursorKind = C.CXCursor_CXXDeleteExpr
	// A unary expression. (noexcept, sizeof, or other traits)
	Cursor_UnaryExpr CursorKind = C.CXCursor_UnaryExpr
	// An Objective-C string literal i.e. @"foo".
	Cursor_ObjCStringLiteral CursorKind = C.CXCursor_ObjCStringLiteral
	// An Objective-C \@encode expression.
	Cursor_ObjCEncodeExpr CursorKind = C.CXCursor_ObjCEncodeExpr
	// An Objective-C \@selector expression.
	Cursor_ObjCSelectorExpr CursorKind = C.CXCursor_ObjCSelectorExpr
	// An Objective-C \@protocol expression.
	Cursor_ObjCProtocolExpr CursorKind = C.CXCursor_ObjCProtocolExpr
	/*
		An Objective-C "bridged" cast expression, which casts between
		Objective-C pointers and C pointers, transferring ownership in the process.

		\code
		NSString *str = (__bridge_transfer NSString *)CFCreateString();
		\endcode
	*/
	Cursor_ObjCBridgedCastExpr CursorKind = C.CXCursor_ObjCBridgedCastExpr
	/*
		Represents a C++0x pack expansion that produces a sequence of
		expressions.

		A pack expansion expression contains a pattern (which itself is an
		expression) followed by an ellipsis. For example:

		\code
		template<typename F, typename ...Types>
		void forward(F f, Types &&...args) {
		f(static_cast<Types&&>(args)...);
		}
		\endcode
	*/
	Cursor_PackExpansionExpr CursorKind = C.CXCursor_PackExpansionExpr
	/*
		Represents an expression that computes the length of a parameter
		pack.

		\code
		template<typename ...Types>
		struct count {
		static const unsigned value = sizeof...(Types);
		};
		\endcode
	*/
	Cursor_SizeOfPackExpr CursorKind = C.CXCursor_SizeOfPackExpr
	Cursor_LambdaExpr     CursorKind = C.CXCursor_LambdaExpr
	// Objective-c Boolean Literal.
	Cursor_ObjCBoolLiteralExpr CursorKind = C.CXCursor_ObjCBoolLiteralExpr
	// Represents the "self" expression in an Objective-C method.
	Cursor_ObjCSelfExpr CursorKind = C.CXCursor_ObjCSelfExpr
	// OpenMP 4.0 [2.4, Array Section].
	Cursor_OMPArraySectionExpr CursorKind = C.CXCursor_OMPArraySectionExpr
	// Represents an @available(...) check.
	Cursor_ObjCAvailabilityCheckExpr CursorKind = C.CXCursor_ObjCAvailabilityCheckExpr
	/*
		A statement whose specific kind is not exposed via this
		interface.

		Unexposed statements have the same operations as any other kind of
		statement; one can extract their location information, spelling,
		children, etc. However, the specific kind of the statement is not
		reported.
	*/
	Cursor_UnexposedStmt CursorKind = C.CXCursor_UnexposedStmt
	/*
		A labelled statement in a function.

		This cursor kind is used to describe the "start_over:" label statement in
		the following example:

		\code
		start_over:
		++counter;
		\endcode
	*/
	Cursor_LabelStmt CursorKind = C.CXCursor_LabelStmt
	/*
		A group of statements like { stmt stmt }.

		This cursor kind is used to describe compound statements, e.g. function
		bodies.
	*/
	Cursor_CompoundStmt CursorKind = C.CXCursor_CompoundStmt
	// A case statement.
	Cursor_CaseStmt CursorKind = C.CXCursor_CaseStmt
	// A default statement.
	Cursor_DefaultStmt CursorKind = C.CXCursor_DefaultStmt
	// An if statement
	Cursor_IfStmt CursorKind = C.CXCursor_IfStmt
	// A switch statement.
	Cursor_SwitchStmt CursorKind = C.CXCursor_SwitchStmt
	// A while statement.
	Cursor_WhileStmt CursorKind = C.CXCursor_WhileStmt
	// A do statement.
	Cursor_DoStmt CursorKind = C.CXCursor_DoStmt
	// A for statement.
	Cursor_ForStmt CursorKind = C.CXCursor_ForStmt
	// A goto statement.
	Cursor_GotoStmt CursorKind = C.CXCursor_GotoStmt
	// An indirect goto statement.
	Cursor_IndirectGotoStmt CursorKind = C.CXCursor_IndirectGotoStmt
	// A continue statement.
	Cursor_ContinueStmt CursorKind = C.CXCursor_ContinueStmt
	// A break statement.
	Cursor_BreakStmt CursorKind = C.CXCursor_BreakStmt
	// A return statement.
	Cursor_ReturnStmt CursorKind = C.CXCursor_ReturnStmt
	// A GCC inline assembly statement extension.
	Cursor_GCCAsmStmt CursorKind = C.CXCursor_GCCAsmStmt
	// A GCC inline assembly statement extension.
	Cursor_AsmStmt CursorKind = C.CXCursor_AsmStmt
	// Objective-C's overall \@try-\@catch-\@finally statement.
	Cursor_ObjCAtTryStmt CursorKind = C.CXCursor_ObjCAtTryStmt
	// Objective-C's \@catch statement.
	Cursor_ObjCAtCatchStmt CursorKind = C.CXCursor_ObjCAtCatchStmt
	// Objective-C's \@finally statement.
	Cursor_ObjCAtFinallyStmt CursorKind = C.CXCursor_ObjCAtFinallyStmt
	// Objective-C's \@throw statement.
	Cursor_ObjCAtThrowStmt CursorKind = C.CXCursor_ObjCAtThrowStmt
	// Objective-C's \@synchronized statement.
	Cursor_ObjCAtSynchronizedStmt CursorKind = C.CXCursor_ObjCAtSynchronizedStmt
	// Objective-C's autorelease pool statement.
	Cursor_ObjCAutoreleasePoolStmt CursorKind = C.CXCursor_ObjCAutoreleasePoolStmt
	// Objective-C's collection statement.
	Cursor_ObjCForCollectionStmt CursorKind = C.CXCursor_ObjCForCollectionStmt
	// C++'s catch statement.
	Cursor_CXXCatchStmt CursorKind = C.CXCursor_CXXCatchStmt
	// C++'s try statement.
	Cursor_CXXTryStmt CursorKind = C.CXCursor_CXXTryStmt
	// C++'s for (* : *) statement.
	Cursor_CXXForRangeStmt CursorKind = C.CXCursor_CXXForRangeStmt
	// Windows Structured Exception Handling's try statement.
	Cursor_SEHTryStmt CursorKind = C.CXCursor_SEHTryStmt
	// Windows Structured Exception Handling's except statement.
	Cursor_SEHExceptStmt CursorKind = C.CXCursor_SEHExceptStmt
	// Windows Structured Exception Handling's finally statement.
	Cursor_SEHFinallyStmt CursorKind = C.CXCursor_SEHFinallyStmt
	// A MS inline assembly statement extension.
	Cursor_MSAsmStmt CursorKind = C.CXCursor_MSAsmStmt
	/*
		The null statement ";": C99 6.8.3p3.

		This cursor kind is used to describe the null statement.
	*/
	Cursor_NullStmt CursorKind = C.CXCursor_NullStmt
	// Adaptor class for mixing declarations with statements and expressions.
	Cursor_DeclStmt CursorKind = C.CXCursor_DeclStmt
	// OpenMP parallel directive.
	Cursor_OMPParallelDirective CursorKind = C.CXCursor_OMPParallelDirective
	// OpenMP SIMD directive.
	Cursor_OMPSimdDirective CursorKind = C.CXCursor_OMPSimdDirective
	// OpenMP for directive.
	Cursor_OMPForDirective CursorKind = C.CXCursor_OMPForDirective
	// OpenMP sections directive.
	Cursor_OMPSectionsDirective CursorKind = C.CXCursor_OMPSectionsDirective
	// OpenMP section directive.
	Cursor_OMPSectionDirective CursorKind = C.CXCursor_OMPSectionDirective
	// OpenMP single directive.
	Cursor_OMPSingleDirective CursorKind = C.CXCursor_OMPSingleDirective
	// OpenMP parallel for directive.
	Cursor_OMPParallelForDirective CursorKind = C.CXCursor_OMPParallelForDirective
	// OpenMP parallel sections directive.
	Cursor_OMPParallelSectionsDirective CursorKind = C.CXCursor_OMPParallelSectionsDirective
	// OpenMP task directive.
	Cursor_OMPTaskDirective CursorKind = C.CXCursor_OMPTaskDirective
	// OpenMP master directive.
	Cursor_OMPMasterDirective CursorKind = C.CXCursor_OMPMasterDirective
	// OpenMP critical directive.
	Cursor_OMPCriticalDirective CursorKind = C.CXCursor_OMPCriticalDirective
	// OpenMP taskyield directive.
	Cursor_OMPTaskyieldDirective CursorKind = C.CXCursor_OMPTaskyieldDirective
	// OpenMP barrier directive.
	Cursor_OMPBarrierDirective CursorKind = C.CXCursor_OMPBarrierDirective
	// OpenMP taskwait directive.
	Cursor_OMPTaskwaitDirective CursorKind = C.CXCursor_OMPTaskwaitDirective
	// OpenMP flush directive.
	Cursor_OMPFlushDirective CursorKind = C.CXCursor_OMPFlushDirective
	// Windows Structured Exception Handling's leave statement.
	Cursor_SEHLeaveStmt CursorKind = C.CXCursor_SEHLeaveStmt
	// OpenMP ordered directive.
	Cursor_OMPOrderedDirective CursorKind = C.CXCursor_OMPOrderedDirective
	// OpenMP atomic directive.
	Cursor_OMPAtomicDirective CursorKind = C.CXCursor_OMPAtomicDirective
	// OpenMP for SIMD directive.
	Cursor_OMPForSimdDirective CursorKind = C.CXCursor_OMPForSimdDirective
	// OpenMP parallel for SIMD directive.
	Cursor_OMPParallelForSimdDirective CursorKind = C.CXCursor_OMPParallelForSimdDirective
	// OpenMP target directive.
	Cursor_OMPTargetDirective CursorKind = C.CXCursor_OMPTargetDirective
	// OpenMP teams directive.
	Cursor_OMPTeamsDirective CursorKind = C.CXCursor_OMPTeamsDirective
	// OpenMP taskgroup directive.
	Cursor_OMPTaskgroupDirective CursorKind = C.CXCursor_OMPTaskgroupDirective
	// OpenMP cancellation point directive.
	Cursor_OMPCancellationPointDirective CursorKind = C.CXCursor_OMPCancellationPointDirective
	// OpenMP cancel directive.
	Cursor_OMPCancelDirective CursorKind = C.CXCursor_OMPCancelDirective
	// OpenMP target data directive.
	Cursor_OMPTargetDataDirective CursorKind = C.CXCursor_OMPTargetDataDirective
	// OpenMP taskloop directive.
	Cursor_OMPTaskLoopDirective CursorKind = C.CXCursor_OMPTaskLoopDirective
	// OpenMP taskloop simd directive.
	Cursor_OMPTaskLoopSimdDirective CursorKind = C.CXCursor_OMPTaskLoopSimdDirective
	// OpenMP distribute directive.
	Cursor_OMPDistributeDirective CursorKind = C.CXCursor_OMPDistributeDirective
	// OpenMP target enter data directive.
	Cursor_OMPTargetEnterDataDirective CursorKind = C.CXCursor_OMPTargetEnterDataDirective
	// OpenMP target exit data directive.
	Cursor_OMPTargetExitDataDirective CursorKind = C.CXCursor_OMPTargetExitDataDirective
	// OpenMP target parallel directive.
	Cursor_OMPTargetParallelDirective CursorKind = C.CXCursor_OMPTargetParallelDirective
	// OpenMP target parallel for directive.
	Cursor_OMPTargetParallelForDirective CursorKind = C.CXCursor_OMPTargetParallelForDirective
	// OpenMP target update directive.
	Cursor_OMPTargetUpdateDirective CursorKind = C.CXCursor_OMPTargetUpdateDirective
	// OpenMP distribute parallel for directive.
	Cursor_OMPDistributeParallelForDirective CursorKind = C.CXCursor_OMPDistributeParallelForDirective
	// OpenMP distribute parallel for simd directive.
	Cursor_OMPDistributeParallelForSimdDirective CursorKind = C.CXCursor_OMPDistributeParallelForSimdDirective
	// OpenMP distribute simd directive.
	Cursor_OMPDistributeSimdDirective CursorKind = C.CXCursor_OMPDistributeSimdDirective
	// OpenMP target parallel for simd directive.
	Cursor_OMPTargetParallelForSimdDirective CursorKind = C.CXCursor_OMPTargetParallelForSimdDirective
	// OpenMP target simd directive.
	Cursor_OMPTargetSimdDirective CursorKind = C.CXCursor_OMPTargetSimdDirective
	// OpenMP teams distribute directive.
	Cursor_OMPTeamsDistributeDirective CursorKind = C.CXCursor_OMPTeamsDistributeDirective
	// OpenMP teams distribute simd directive.
	Cursor_OMPTeamsDistributeSimdDirective CursorKind = C.CXCursor_OMPTeamsDistributeSimdDirective
	// OpenMP teams distribute parallel for simd directive.
	Cursor_OMPTeamsDistributeParallelForSimdDirective CursorKind = C.CXCursor_OMPTeamsDistributeParallelForSimdDirective
	// OpenMP teams distribute parallel for directive.
	Cursor_OMPTeamsDistributeParallelForDirective CursorKind = C.CXCursor_OMPTeamsDistributeParallelForDirective
	// OpenMP target teams directive.
	Cursor_OMPTargetTeamsDirective CursorKind = C.CXCursor_OMPTargetTeamsDirective
	// OpenMP target teams distribute directive.
	Cursor_OMPTargetTeamsDistributeDirective CursorKind = C.CXCursor_OMPTargetTeamsDistributeDirective
	// OpenMP target teams distribute parallel for directive.
	Cursor_OMPTargetTeamsDistributeParallelForDirective CursorKind = C.CXCursor_OMPTargetTeamsDistributeParallelForDirective
	// OpenMP target teams distribute parallel for simd directive.
	Cursor_OMPTargetTeamsDistributeParallelForSimdDirective CursorKind = C.CXCursor_OMPTargetTeamsDistributeParallelForSimdDirective
	// OpenMP target teams distribute simd directive.
	Cursor_OMPTargetTeamsDistributeSimdDirective CursorKind = C.CXCursor_OMPTargetTeamsDistributeSimdDirective
	/*
		Cursor that represents the translation unit itself.

		The translation unit cursor exists primarily to act as the root
		cursor for traversing the contents of a translation unit.
	*/
	Cursor_TranslationUnit CursorKind = C.CXCursor_TranslationUnit
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_UnexposedAttr CursorKind = C.CXCursor_UnexposedAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_IBActionAttr CursorKind = C.CXCursor_IBActionAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_IBOutletAttr CursorKind = C.CXCursor_IBOutletAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_IBOutletCollectionAttr CursorKind = C.CXCursor_IBOutletCollectionAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_CXXFinalAttr CursorKind = C.CXCursor_CXXFinalAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_CXXOverrideAttr CursorKind = C.CXCursor_CXXOverrideAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_AnnotateAttr CursorKind = C.CXCursor_AnnotateAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_AsmLabelAttr CursorKind = C.CXCursor_AsmLabelAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_PackedAttr CursorKind = C.CXCursor_PackedAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_PureAttr CursorKind = C.CXCursor_PureAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_ConstAttr CursorKind = C.CXCursor_ConstAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_NoDuplicateAttr CursorKind = C.CXCursor_NoDuplicateAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_CUDAConstantAttr CursorKind = C.CXCursor_CUDAConstantAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_CUDADeviceAttr CursorKind = C.CXCursor_CUDADeviceAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_CUDAGlobalAttr CursorKind = C.CXCursor_CUDAGlobalAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_CUDAHostAttr CursorKind = C.CXCursor_CUDAHostAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_CUDASharedAttr CursorKind = C.CXCursor_CUDASharedAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_VisibilityAttr CursorKind = C.CXCursor_VisibilityAttr
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_DLLExport CursorKind = C.CXCursor_DLLExport
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_DLLImport CursorKind = C.CXCursor_DLLImport
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_PreprocessingDirective CursorKind = C.CXCursor_PreprocessingDirective
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_MacroDefinition CursorKind = C.CXCursor_MacroDefinition
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_MacroExpansion CursorKind = C.CXCursor_MacroExpansion
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_MacroInstantiation CursorKind = C.CXCursor_MacroInstantiation
	// An attribute whose specific kind is not exposed via this interface.
	Cursor_InclusionDirective CursorKind = C.CXCursor_InclusionDirective
	// A module import declaration.
	Cursor_ModuleImportDecl CursorKind = C.CXCursor_ModuleImportDecl
	// A module import declaration.
	Cursor_TypeAliasTemplateDecl CursorKind = C.CXCursor_TypeAliasTemplateDecl
	// A static_assert or _Static_assert node
	Cursor_StaticAssert CursorKind = C.CXCursor_StaticAssert
	// a friend declaration.
	Cursor_FriendDecl CursorKind = C.CXCursor_FriendDecl
	// A code completion overload candidate.
	Cursor_OverloadCandidate CursorKind = C.CXCursor_OverloadCandidate

	// Extras we don't need strings for. Primarily because they are place holders
	// and have values that conflict with actual constants that we want to track.
	Cursor_FirstDecl          CursorKind = C.CXCursor_FirstDecl
	Cursor_LastDecl           CursorKind = C.CXCursor_LastDecl
	Cursor_FirstRef           CursorKind = C.CXCursor_FirstRef
	Cursor_LastRef            CursorKind = C.CXCursor_LastRef
	Cursor_FirstInvalid       CursorKind = C.CXCursor_FirstInvalid
	Cursor_LastInvalid        CursorKind = C.CXCursor_LastInvalid
	Cursor_FirstExpr          CursorKind = C.CXCursor_FirstExpr
	Cursor_LastExpr           CursorKind = C.CXCursor_LastExpr
	Cursor_FirstStmt          CursorKind = C.CXCursor_FirstStmt
	Cursor_LastStmt           CursorKind = C.CXCursor_LastStmt
	Cursor_FirstAttr          CursorKind = C.CXCursor_FirstAttr
	Cursor_LastAttr           CursorKind = C.CXCursor_LastAttr
	Cursor_FirstPreprocessing CursorKind = C.CXCursor_FirstPreprocessing
	Cursor_LastPreprocessing  CursorKind = C.CXCursor_LastPreprocessing
	Cursor_FirstExtraDecl     CursorKind = C.CXCursor_FirstExtraDecl
	Cursor_LastExtraDecl      CursorKind = C.CXCursor_LastExtraDecl
)

// Determine whether the given cursor kind represents a normal declaration,
// not an extra declaration.
func (ck CursorKind) IsNormDeclaration() bool {
	return Cursor_FirstDecl <= ck && ck <= Cursor_LastDecl
}

// Determine whether the given cursor kind represents an extra declaration.
func (ck CursorKind) IsExtraDeclaration() bool {
	return Cursor_FirstExtraDecl <= ck && ck <= Cursor_LastExtraDecl
}

// Determine whether the given cursor kind represents a declaration.
func (ck CursorKind) IsDeclaration() bool {
	return ck.IsNormDeclaration() || ck.IsExtraDeclaration()
}

func (ck CursorKind) IsDeclarationSlow() bool {
	o := C.clang_isDeclaration(C.enum_CXCursorKind(ck))

	return o != C.uint(0)
}

/*
	Determine whether the given cursor kind represents a simple
	reference.

	Note that other kinds of cursors (such as expressions) can also refer to
	other cursors. Use clang_getCursorReferenced() to determine whether a
	particular cursor refers to another entity.
*/
func (ck CursorKind) IsReference() bool {
	return Cursor_FirstRef <= ck && ck <= Cursor_LastRef
}

func (ck CursorKind) IsReferenceSlow() bool {
	o := C.clang_isReference(C.enum_CXCursorKind(ck))

	return o != C.uint(0)
}

// Determine whether the given cursor kind represents an expression.
func (ck CursorKind) IsExpression() bool {
	return Cursor_FirstExpr <= ck && ck <= Cursor_LastExpr
}

func (ck CursorKind) IsExpressionSlow() bool {
	o := C.clang_isExpression(C.enum_CXCursorKind(ck))

	return o != C.uint(0)
}

// Determine whether the given cursor kind represents a statement.
func (ck CursorKind) IsStatement() bool {
	return Cursor_FirstStmt <= ck && ck <= Cursor_LastStmt
}

func (ck CursorKind) IsStatementSlow() bool {
	o := C.clang_isStatement(C.enum_CXCursorKind(ck))

	return o != C.uint(0)
}

// Determine whether the given cursor kind represents an attribute.
func (ck CursorKind) IsAttribute() bool {
	return Cursor_FirstAttr <= ck && ck <= Cursor_LastAttr
}

func (ck CursorKind) IsAttributeSlow() bool {
	o := C.clang_isAttribute(C.enum_CXCursorKind(ck))

	return o != C.uint(0)
}

// Determine whether the given cursor kind represents an invalid cursor.
func (ck CursorKind) IsInvalid() bool {
	return Cursor_FirstInvalid <= ck && ck <= Cursor_LastInvalid
}

func (ck CursorKind) IsInvalidSlow() bool {
	o := C.clang_isInvalid(C.enum_CXCursorKind(ck))

	return o != C.uint(0)
}

// Determine whether the given cursor kind represents a translation unit.
func (ck CursorKind) IsTranslationUnit() bool {
	return ck == Cursor_TranslationUnit
}

func (ck CursorKind) IsTranslationUnitSlow() bool {
	o := C.clang_isTranslationUnit(C.enum_CXCursorKind(ck))

	return o != C.uint(0)
}

// * Determine whether the given cursor represents a preprocessing element, such as a preprocessor directive or macro instantiation.
func (ck CursorKind) IsPreprocessing() bool {
	return Cursor_FirstPreprocessing <= ck && ck <= Cursor_LastPreprocessing
}

func (ck CursorKind) IsPreprocessingSlow() bool {
	o := C.clang_isPreprocessing(C.enum_CXCursorKind(ck))

	return o != C.uint(0)
}

// * Determine whether the given cursor represents a currently unexposed piece of the AST (e.g., CXCursor_UnexposedStmt).
func (ck CursorKind) IsUnexposed() bool {
	switch ck {
	case Cursor_UnexposedDecl, Cursor_UnexposedExpr, Cursor_UnexposedStmt, Cursor_UnexposedAttr:
		return true
	}
	return false
}

func (ck CursorKind) IsUnexposedSlow() bool {
	o := C.clang_isUnexposed(C.enum_CXCursorKind(ck))

	return o != C.uint(0)
}

// IsLiteral returns true for Literal kinds; these return no Spelling string, but getting the token is an option.
// There are some CursorKinds with the name Literal in them that sound more like expressions so there weren't
// included here.
func (ck CursorKind) IsLiteral() bool {
	switch ck {
	case Cursor_IntegerLiteral,
		Cursor_FloatingLiteral,
		Cursor_ImaginaryLiteral,
		Cursor_StringLiteral,
		Cursor_CharacterLiteral,
		Cursor_ObjCStringLiteral:
		return true
	}
	return false
}

func (ck CursorKind) Spelling() string {
	return cx2GoString(C.clang_getCursorKindSpelling(C.enum_CXCursorKind(ck)))
}
