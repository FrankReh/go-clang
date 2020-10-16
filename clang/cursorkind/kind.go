package cursorkind

type Kind int

//go:generate stringer -type=Kind

//
// Describes the kind of entity that a cursor refers to.
//
const (
	/* Declarations */
	/**
	 * A declaration whose specific kind is not exposed via this
	 * interface.
	 *
	 * Unexposed declarations have the same operations as any other kind
	 * of declaration; one can extract their location information
	 * spelling, find their definitions, etc. However, the specific kind
	 * of the declaration is not reported.
	 */
	UnexposedDecl Kind = 1
	StructDecl    Kind = 2 /** A C or C++ struct. */
	UnionDecl     Kind = 3 /** A C or C++ union. */
	ClassDecl     Kind = 4 /** A C++ class. */
	EnumDecl      Kind = 5 /** An enumeration. */
	/**
	 * A field (in C) or non-static data member (in C++) in a
	 * struct, union, or C++ class.
	 */
	FieldDecl                          Kind = 6
	EnumConstantDecl                   Kind = 7  /** An enumerator constant. */
	FunctionDecl                       Kind = 8  /** A function. */
	VarDecl                            Kind = 9  /** A variable. */
	ParmDecl                           Kind = 10 /** A function or method parameter. */
	ObjCInterfaceDecl                  Kind = 11 /** An Objective-C \@interface. */
	ObjCCategoryDecl                   Kind = 12 /** An Objective-C \@interface for a category. */
	ObjCProtocolDecl                   Kind = 13 /** An Objective-C \@protocol declaration. */
	ObjCPropertyDecl                   Kind = 14 /** An Objective-C \@property declaration. */
	ObjCIvarDecl                       Kind = 15 /** An Objective-C instance variable. */
	ObjCInstanceMethodDecl             Kind = 16 /** An Objective-C instance method. */
	ObjCClassMethodDecl                Kind = 17 /** An Objective-C class method. */
	ObjCImplementationDecl             Kind = 18 /** An Objective-C \@implementation. */
	ObjCCategoryImplDecl               Kind = 19 /** An Objective-C \@implementation for a category. */
	TypedefDecl                        Kind = 20 /** A typedef. */
	CXXMethod                          Kind = 21 /** A C++ class method. */
	Namespace                          Kind = 22 /** A C++ namespace. */
	LinkageSpec                        Kind = 23 /** A linkage specification, e.g. 'extern "C"'. */
	Constructor                        Kind = 24 /** A C++ constructor. */
	Destructor                         Kind = 25 /** A C++ destructor. */
	ConversionFunction                 Kind = 26 /** A C++ conversion function. */
	TemplateTypeParameter              Kind = 27 /** A C++ template type parameter. */
	NonTypeTemplateParameter           Kind = 28 /** A C++ non-type template parameter. */
	TemplateTemplateParameter          Kind = 29 /** A C++ template template parameter. */
	FunctionTemplate                   Kind = 30 /** A C++ function template. */
	ClassTemplate                      Kind = 31 /** A C++ class template. */
	ClassTemplatePartialSpecialization Kind = 32 /** A C++ class template partial specialization. */
	NamespaceAlias                     Kind = 33 /** A C++ namespace alias declaration. */
	UsingDirective                     Kind = 34 /** A C++ using directive. */
	UsingDeclaration                   Kind = 35 /** A C++ using declaration. */
	TypeAliasDecl                      Kind = 36 /** A C++ alias declaration */
	ObjCSynthesizeDecl                 Kind = 37 /** An Objective-C \@synthesize definition. */
	ObjCDynamicDecl                    Kind = 38 /** An Objective-C \@dynamic definition. */
	CXXAccessSpecifier                 Kind = 39 /** An access specifier. */

	/* References */
	ObjCSuperClassRef Kind = 40 /* Decl references */
	ObjCProtocolRef   Kind = 41
	ObjCClassRef      Kind = 42
	/**
	 * A reference to a type declaration.
	 *
	 * A type reference occurs anywhere where a type is named but not
	 * declared. For example, given:
	 *
	 * \code
	 * typedef unsigned size_type;
	 * size_type size;
	 * \endcode
	 *
	 * The typedef is a declaration of size_type (TypedefDecl)
	 * while the type of the variable "size" is referenced. The cursor
	 * referenced by the type of size is the typedef for size_type.
	 */
	TypeRef          Kind = 43
	CXXBaseSpecifier Kind = 44
	/**
	 * A reference to a class template, function template, template
	 * template parameter, or class template partial specialization.
	 */
	TemplateRef Kind = 45
	/**
	 * A reference to a namespace or namespace alias.
	 */
	NamespaceRef Kind = 46
	/**
	 * A reference to a member of a struct, union, or class that occurs in
	 * some non-expression context, e.g., a designated initializer.
	 */
	MemberRef Kind = 47
	/**
	 * A reference to a labeled statement.
	 *
	 * This cursor kind is used to describe the jump to "start_over" in the
	 * goto statement in the following example:
	 *
	 * \code
	 *   start_over:
	 *     ++counter;
	 *
	 *     goto start_over;
	 * \endcode
	 *
	 * A label reference cursor refers to a label statement.
	 */
	LabelRef Kind = 48

	/**
	 * A reference to a set of overloaded functions or function templates
	 * that has not yet been resolved to a specific function or function template.
	 *
	 * An overloaded declaration reference cursor occurs in C++ templates where
	 * a dependent name refers to a function. For example:
	 *
	 * \code
	 * template<typename T> void swap(T&, T&);
	 *
	 * struct X { ... };
	 * void swap(X&, X&);
	 *
	 * template<typename T>
	 * void reverse(T* first, T* last) {
	 *   while (first < last - 1) {
	 *     swap(*first, *--last);
	 *     ++first;
	 *   }
	 * }
	 *
	 * struct Y { };
	 * void swap(Y&, Y&);
	 * \endcode
	 *
	 * Here, the identifier "swap" is associated with an overloaded declaration
	 * reference. In the template definition, "swap" refers to either of the two
	 * "swap" functions declared above, so both results will be available. At
	 * instantiation time, "swap" may also refer to other functions found via
	 * argument-dependent lookup (e.g., the "swap" function at the end of the
	 * example).
	 *
	 * The functions \c clang_getNumOverloadedDecls() and
	 * \c clang_getOverloadedDecl() can be used to retrieve the definitions
	 * referenced by this cursor.
	 */
	OverloadedDeclRef Kind = 49

	/**
	 * A reference to a variable that occurs in some non-expression
	 * context, e.g., a C++ lambda capture list.
	 */
	VariableRef Kind = 50

	/* Error conditions */
	InvalidFile    Kind = 70
	NoDeclFound    Kind = 71
	NotImplemented Kind = 72
	InvalidCode    Kind = 73

	/* Expressions */

	/**
	 * An expression whose specific kind is not exposed via this
	 * interface.
	 *
	 * Unexposed expressions have the same operations as any other kind
	 * of expression; one can extract their location information
	 * spelling, children, etc. However, the specific kind of the
	 * expression is not reported.
	 */
	UnexposedExpr Kind = 100

	/**
	 * An expression that refers to some value declaration, such
	 * as a function, variable, or enumerator.
	 */
	DeclRefExpr Kind = 101

	/**
	 * An expression that refers to a member of a struct, union
	 * class, Objective-C class, etc.
	 */
	MemberRefExpr Kind = 102

	/** An expression that calls a function. */
	CallExpr Kind = 103

	/** An expression that sends a message to an Objective-C
	  object or class. */
	ObjCMessageExpr Kind = 104

	/** An expression that represents a block literal. */
	BlockExpr Kind = 105

	/** An integer literal.
	 */
	IntegerLiteral Kind = 106

	/** A floating point number literal.
	 */
	FloatingLiteral Kind = 107

	/** An imaginary number literal.
	 */
	ImaginaryLiteral Kind = 108

	/** A string literal.
	 */
	StringLiteral Kind = 109

	/** A character literal.
	 */
	CharacterLiteral Kind = 110

	/** A parenthesized expression, e.g. "(1)".
	 *
	 * This AST node is only formed if full location information is requested.
	 */
	ParenExpr Kind = 111

	/** This represents the unary-expression's (except sizeof and
	 * alignof).
	 */
	UnaryOperator Kind = 112

	/** [C99 6.5.2.1] Array Subscripting.
	 */
	ArraySubscriptExpr Kind = 113

	/** A builtin binary operation expression such as "x + y" or
	 * "x <Kind = y".
	 */
	BinaryOperator Kind = 114

	/** Compound assignment such as "+Kind =".
	 */
	CompoundAssignOperator Kind = 115

	/** The ?: ternary operator.
	 */
	ConditionalOperator Kind = 116

	/** An explicit cast in C (C99 6.5.4) or a C-style cast in C++
	 * (C++ [expr.cast]), which uses the syntax (Type)expr.
	 *
	 * For example: (int)f.
	 */
	CStyleCastExpr Kind = 117

	/** [C99 6.5.2.5]
	 */
	CompoundLiteralExpr Kind = 118

	/** Describes an C or C++ initializer list.
	 */
	InitListExpr Kind = 119

	/** The GNU address of label extension, representing &&label.
	 */
	AddrLabelExpr Kind = 120

	/** This is the GNU Statement Expression extension: ({int XKind =4; X;})
	 */
	StmtExpr Kind = 121

	/** Represents a C11 generic selection.
	 */
	GenericSelectionExpr Kind = 122

	/** Implements the GNU __null extension, which is a name for a null
	 * pointer constant that has integral type (e.g., int or long) and is the same
	 * size and alignment as a pointer.
	 *
	 * The __null extension is typically only used by system headers, which define
	 * NULL as __null in C++ rather than using 0 (which is an integer that may not
	 * match the size of a pointer).
	 */
	GNUNullExpr Kind = 123

	/** C++'s static_cast<> expression.
	 */
	CXXStaticCastExpr Kind = 124

	/** C++'s dynamic_cast<> expression.
	 */
	CXXDynamicCastExpr Kind = 125

	/** C++'s reinterpret_cast<> expression.
	 */
	CXXReinterpretCastExpr Kind = 126

	/** C++'s const_cast<> expression.
	 */
	CXXConstCastExpr Kind = 127

	/** Represents an explicit C++ type conversion that uses "functional"
	 * notion (C++ [expr.type.conv]).
	 *
	 * Example:
	 * \code
	 *   x = int(0.5);
	 * \endcode
	 */
	CXXFunctionalCastExpr Kind = 128

	/** A C++ typeid expression (C++ [expr.typeid]).
	 */
	CXXTypeidExpr Kind = 129

	/** [C++ 2.13.5] C++ Boolean Literal.
	 */
	CXXBoolLiteralExpr Kind = 130

	/** [C++0x 2.14.7] C++ Pointer Literal.
	 */
	CXXNullPtrLiteralExpr Kind = 131

	/** Represents the "this" expression in C++
	 */
	CXXThisExpr Kind = 132

	/** [C++ 15] C++ Throw Expression.
	 *
	 * This handles 'throw' and 'throw' assignment-expression. When
	 * assignment-expression isn't present, Op will be null.
	 */
	CXXThrowExpr Kind = 133

	/** A new expression for memory allocation and constructor calls, e.g:
	 * "new CXXNewExpr(foo)".
	 */
	CXXNewExpr Kind = 134

	/** A delete expression for memory deallocation and destructor calls
	 * e.g. "delete[] pArray".
	 */
	CXXDeleteExpr Kind = 135

	/** A unary expression. (noexcept, sizeof, or other traits)
	 */
	UnaryExpr Kind = 136

	/** An Objective-C string literal i.e. @"foo".
	 */
	ObjCStringLiteral Kind = 137

	/** An Objective-C \@encode expression.
	 */
	ObjCEncodeExpr Kind = 138

	/** An Objective-C \@selector expression.
	 */
	ObjCSelectorExpr Kind = 139

	/** An Objective-C \@protocol expression.
	 */
	ObjCProtocolExpr Kind = 140

	/** An Objective-C "bridged" cast expression, which casts between
	 * Objective-C pointers and C pointers, transferring ownership in the process.
	 *
	 * \code
	 *   NSString *str Kind = (__bridge_transfer NSString *)CFCreateString();
	 * \endcode
	 */
	ObjCBridgedCastExpr Kind = 141

	/** Represents a C++0x pack expansion that produces a sequence of
	 * expressions.
	 *
	 * A pack expansion expression contains a pattern (which itself is an
	 * expression) followed by an ellipsis. For example:
	 *
	 * \code
	 * template<typename F, typename ...Types>
	 * void forward(F f, Types &&...args) {
	 *  f(static_cast<Types&&>(args)...);
	 * }
	 * \endcode
	 */
	PackExpansionExpr Kind = 142

	/** Represents an expression that computes the length of a parameter
	 * pack.
	 *
	 * \code
	 * template<typename ...Types>
	 * struct count {
	 *   static const unsigned value Kind = sizeof...(Types);
	 * };
	 * \endcode
	 */
	SizeOfPackExpr Kind = 143

	/* Represents a C++ lambda expression that produces a local function
	 * object.
	 *
	 * \code
	 * void abssort(float *x, unsigned N) {
	 *   std::sort(x, x + N
	 *             [](float a, float b) {
	 *               return std::abs(a) < std::abs(b);
	 *             });
	 * }
	 * \endcode
	 */
	LambdaExpr Kind = 144

	/** Objective-c Boolean Literal.
	 */
	ObjCBoolLiteralExpr Kind = 145

	/** Represents the "self" expression in an Objective-C method.
	 */
	ObjCSelfExpr Kind = 146

	/** OpenMP 4.0 [2.4, Array Section].
	 */
	OMPArraySectionExpr Kind = 147

	/** Represents an @available(...) check.
	 */
	ObjCAvailabilityCheckExpr Kind = 148

	/** Fixed point literal
	 */
	FixedPointLiteral = 149

	/* Statements */
	/**
	 * A statement whose specific kind is not exposed via this
	 * interface.
	 *
	 * Unexposed statements have the same operations as any other kind of
	 * statement; one can extract their location information, spelling
	 * children, etc. However, the specific kind of the statement is not
	 * reported.
	 */
	UnexposedStmt Kind = 200

	/** A labelled statement in a function.
	 *
	 * This cursor kind is used to describe the "start_over:" label statement in
	 * the following example:
	 *
	 * \code
	 *   start_over:
	 *     ++counter;
	 * \endcode
	 *
	 */
	LabelStmt Kind = 201

	/** A group of statements like { stmt stmt }.
	 *
	 * This cursor kind is used to describe compound statements, e.g. function
	 * bodies.
	 */
	CompoundStmt Kind = 202

	/** A case statement.
	 */
	CaseStmt Kind = 203

	/** A default statement.
	 */
	DefaultStmt Kind = 204

	/** An if statement
	 */
	IfStmt Kind = 205

	/** A switch statement.
	 */
	SwitchStmt Kind = 206

	/** A while statement.
	 */
	WhileStmt Kind = 207

	/** A do statement.
	 */
	DoStmt Kind = 208

	/** A for statement.
	 */
	ForStmt Kind = 209

	/** A goto statement.
	 */
	GotoStmt Kind = 210

	/** An indirect goto statement.
	 */
	IndirectGotoStmt Kind = 211

	/** A continue statement.
	 */
	ContinueStmt Kind = 212

	/** A break statement.
	 */
	BreakStmt Kind = 213

	/** A return statement.
	 */
	ReturnStmt Kind = 214

	/** A GCC inline assembly statement extension.
	 */
	GCCAsmStmt Kind = 215
	AsmStmt    Kind = GCCAsmStmt

	/** Objective-C's overall \@try-\@catch-\@finally statement.
	 */
	ObjCAtTryStmt Kind = 216

	/** Objective-C's \@catch statement.
	 */
	ObjCAtCatchStmt Kind = 217

	/** Objective-C's \@finally statement.
	 */
	ObjCAtFinallyStmt Kind = 218

	/** Objective-C's \@throw statement.
	 */
	ObjCAtThrowStmt Kind = 219

	/** Objective-C's \@synchronized statement.
	 */
	ObjCAtSynchronizedStmt Kind = 220

	/** Objective-C's autorelease pool statement.
	 */
	ObjCAutoreleasePoolStmt Kind = 221

	/** Objective-C's collection statement.
	 */
	ObjCForCollectionStmt Kind = 222

	/** C++'s catch statement.
	 */
	CXXCatchStmt Kind = 223

	/** C++'s try statement.
	 */
	CXXTryStmt Kind = 224

	/** C++'s for (* : *) statement.
	 */
	CXXForRangeStmt Kind = 225

	/** Windows Structured Exception Handling's try statement.
	 */
	SEHTryStmt Kind = 226

	/** Windows Structured Exception Handling's except statement.
	 */
	SEHExceptStmt Kind = 227

	/** Windows Structured Exception Handling's finally statement.
	 */
	SEHFinallyStmt Kind = 228

	/** A MS inline assembly statement extension.
	 */
	MSAsmStmt Kind = 229

	/** The null statement ";": C99 6.8.3p3.
	 *
	 * This cursor kind is used to describe the null statement.
	 */
	NullStmt Kind = 230

	/** Adaptor class for mixing declarations with statements and
	 * expressions.
	 */
	DeclStmt Kind = 231

	/** OpenMP parallel directive.
	 */
	OMPParallelDirective Kind = 232

	/** OpenMP SIMD directive.
	 */
	OMPSimdDirective Kind = 233

	/** OpenMP for directive.
	 */
	OMPForDirective Kind = 234

	/** OpenMP sections directive.
	 */
	OMPSectionsDirective Kind = 235

	/** OpenMP section directive.
	 */
	OMPSectionDirective Kind = 236

	/** OpenMP single directive.
	 */
	OMPSingleDirective Kind = 237

	/** OpenMP parallel for directive.
	 */
	OMPParallelForDirective Kind = 238

	/** OpenMP parallel sections directive.
	 */
	OMPParallelSectionsDirective Kind = 239

	/** OpenMP task directive.
	 */
	OMPTaskDirective Kind = 240

	/** OpenMP master directive.
	 */
	OMPMasterDirective Kind = 241

	/** OpenMP critical directive.
	 */
	OMPCriticalDirective Kind = 242

	/** OpenMP taskyield directive.
	 */
	OMPTaskyieldDirective Kind = 243

	/** OpenMP barrier directive.
	 */
	OMPBarrierDirective Kind = 244

	/** OpenMP taskwait directive.
	 */
	OMPTaskwaitDirective Kind = 245

	/** OpenMP flush directive.
	 */
	OMPFlushDirective Kind = 246

	/** Windows Structured Exception Handling's leave statement.
	 */
	SEHLeaveStmt Kind = 247

	/** OpenMP ordered directive.
	 */
	OMPOrderedDirective Kind = 248

	/** OpenMP atomic directive.
	 */
	OMPAtomicDirective Kind = 249

	/** OpenMP for SIMD directive.
	 */
	OMPForSimdDirective Kind = 250

	/** OpenMP parallel for SIMD directive.
	 */
	OMPParallelForSimdDirective Kind = 251

	/** OpenMP target directive.
	 */
	OMPTargetDirective Kind = 252

	/** OpenMP teams directive.
	 */
	OMPTeamsDirective Kind = 253

	/** OpenMP taskgroup directive.
	 */
	OMPTaskgroupDirective Kind = 254

	/** OpenMP cancellation point directive.
	 */
	OMPCancellationPointDirective Kind = 255

	/** OpenMP cancel directive.
	 */
	OMPCancelDirective Kind = 256

	/** OpenMP target data directive.
	 */
	OMPTargetDataDirective Kind = 257

	/** OpenMP taskloop directive.
	 */
	OMPTaskLoopDirective Kind = 258

	/** OpenMP taskloop simd directive.
	 */
	OMPTaskLoopSimdDirective Kind = 259

	/** OpenMP distribute directive.
	 */
	OMPDistributeDirective Kind = 260

	/** OpenMP target enter data directive.
	 */
	OMPTargetEnterDataDirective Kind = 261

	/** OpenMP target exit data directive.
	 */
	OMPTargetExitDataDirective Kind = 262

	/** OpenMP target parallel directive.
	 */
	OMPTargetParallelDirective Kind = 263

	/** OpenMP target parallel for directive.
	 */
	OMPTargetParallelForDirective Kind = 264

	/** OpenMP target update directive.
	 */
	OMPTargetUpdateDirective Kind = 265

	/** OpenMP distribute parallel for directive.
	 */
	OMPDistributeParallelForDirective Kind = 266

	/** OpenMP distribute parallel for simd directive.
	 */
	OMPDistributeParallelForSimdDirective Kind = 267

	/** OpenMP distribute simd directive.
	 */
	OMPDistributeSimdDirective Kind = 268

	/** OpenMP target parallel for simd directive.
	 */
	OMPTargetParallelForSimdDirective Kind = 269

	/** OpenMP target simd directive.
	 */
	OMPTargetSimdDirective Kind = 270

	/** OpenMP teams distribute directive.
	 */
	OMPTeamsDistributeDirective Kind = 271

	/** OpenMP teams distribute simd directive.
	 */
	OMPTeamsDistributeSimdDirective Kind = 272

	/** OpenMP teams distribute parallel for simd directive.
	 */
	OMPTeamsDistributeParallelForSimdDirective Kind = 273

	/** OpenMP teams distribute parallel for directive.
	 */
	OMPTeamsDistributeParallelForDirective Kind = 274

	/** OpenMP target teams directive.
	 */
	OMPTargetTeamsDirective Kind = 275

	/** OpenMP target teams distribute directive.
	 */
	OMPTargetTeamsDistributeDirective Kind = 276

	/** OpenMP target teams distribute parallel for directive.
	 */
	OMPTargetTeamsDistributeParallelForDirective Kind = 277

	/** OpenMP target teams distribute parallel for simd directive.
	 */
	OMPTargetTeamsDistributeParallelForSimdDirective Kind = 278

	/** OpenMP target teams distribute simd directive.
	 */
	OMPTargetTeamsDistributeSimdDirective Kind = 279

	/** C++2a std::bit_cast expression.
	 */
	BuiltinBitCastExpr Kind = 280

	/**
	 * Cursor that represents the translation unit itself.
	 *
	 * The translation unit cursor exists primarily to act as the root
	 * cursor for traversing the contents of a translation unit.
	 */
	TranslationUnit Kind = 300

	/* Attributes */
	/**
	 * An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	UnexposedAttr Kind = 400

	IBActionAttr              Kind = 401
	IBOutletAttr              Kind = 402
	IBOutletCollectionAttr    Kind = 403
	CXXFinalAttr              Kind = 404
	CXXOverrideAttr           Kind = 405
	AnnotateAttr              Kind = 406
	AsmLabelAttr              Kind = 407
	PackedAttr                Kind = 408
	PureAttr                  Kind = 409
	ConstAttr                 Kind = 410
	NoDuplicateAttr           Kind = 411
	CUDAConstantAttr          Kind = 412
	CUDADeviceAttr            Kind = 413
	CUDAGlobalAttr            Kind = 414
	CUDAHostAttr              Kind = 415
	CUDASharedAttr            Kind = 416
	VisibilityAttr            Kind = 417
	DLLExport                 Kind = 418
	DLLImport                 Kind = 419
	NSReturnsRetained         Kind = 420
	NSReturnsNotRetained      Kind = 421
	NSReturnsAutoreleased     Kind = 422
	NSConsumesSelf            Kind = 423
	NSConsumed                Kind = 424
	ObjCException             Kind = 425
	ObjCNSObject              Kind = 426
	ObjCIndependentClass      Kind = 427
	ObjCPreciseLifetime       Kind = 428
	ObjCReturnsInnerPointer   Kind = 429
	ObjCRequiresSuper         Kind = 430
	ObjCRootClass             Kind = 431
	ObjCSubclassingRestricted Kind = 432
	ObjCExplicitProtocolImpl  Kind = 433
	ObjCDesignatedInitializer Kind = 434
	ObjCRuntimeVisible        Kind = 435
	ObjCBoxable               Kind = 436
	FlagEnum                  Kind = 437
	ConvergentAttr            Kind = 438
	WarnUnusedAttr            Kind = 439
	WarnUnusedResultAttr      Kind = 440
	AlignedAttr               Kind = 441

	/* Preprocessing */
	PreprocessingDirective Kind = 500
	MacroDefinition        Kind = 501
	MacroExpansion         Kind = 502
	MacroInstantiation     Kind = MacroExpansion
	InclusionDirective     Kind = 503

	/* Extra Declarations */
	/**
	 * A module import declaration.
	 */
	ModuleImportDecl      Kind = 600
	TypeAliasTemplateDecl Kind = 601
	/**
	 * A static_assert or _Static_assert node
	 */
	StaticAssert Kind = 602
	/**
	 * a friend declaration.
	 */
	FriendDecl Kind = 603

	/**
	 * A code completion overload candidate.
	 */
	OverloadCandidate Kind = 700

	/**
	 * A go-clang value, created to represent a cursor that points back to a
	 * previous cursor that has already been seen in the recursive walk.
	 * This expressly does not come from the libclang header file.
	 */
	Back Kind = -1

	FirstDecl          Kind = UnexposedDecl
	LastDecl           Kind = CXXAccessSpecifier
	FirstRef           Kind = ObjCSuperClassRef
	LastRef            Kind = VariableRef
	FirstInvalid       Kind = InvalidFile
	LastInvalid        Kind = InvalidCode
	FirstExpr          Kind = UnexposedExpr
	LastExpr           Kind = FixedPointLiteral
	FirstStmt          Kind = UnexposedStmt
	LastStmt           Kind = BuiltinBitCastExpr
	FirstAttr          Kind = UnexposedAttr
	LastAttr           Kind = AlignedAttr
	FirstPreprocessing Kind = PreprocessingDirective
	LastPreprocessing  Kind = InclusionDirective
	FirstExtraDecl     Kind = ModuleImportDecl
	LastExtraDecl      Kind = FriendDecl
)

func Validate(i int) (Kind, error) {
	switch {
	case -1 == i:
		return Kind(i), nil
	case int(FirstDecl) <= i && i <= int(LastRef):
		return Kind(i), nil
	case int(FirstInvalid) <= i && i <= int(LastInvalid):
		return Kind(i), nil
	case int(FirstExpr) <= i && i <= int(LastExpr):
		return Kind(i), nil
	case int(FirstStmt) <= i && i <= int(LastStmt):
		return Kind(i), nil
	case i == int(TranslationUnit):
		return Kind(i), nil
	case int(FirstAttr) <= i && i <= int(LastAttr):
		return Kind(i), nil
	case int(FirstPreprocessing) <= i && i <= int(LastPreprocessing):
		return Kind(i), nil
	case int(FirstExtraDecl) <= i && i <= int(LastExtraDecl):
		return Kind(i), nil
	case i == int(OverloadCandidate):
		return Kind(i), nil
	default:
		return 0, InvalidErr
	}
}

func MustValidate(i int) Kind {
	kind, err := Validate(i)
	if err != nil {
		panic(err.Error() + " " + Kind(i).String())
	}
	return kind
}

type Error string

const InvalidErr Error = "Invalid value"

func (e Error) Error() string {
	return string(e)
}

// Determine whether the given cursor kind represents a normal declaration,
// not an extra declaration.
func (ck Kind) IsNormDeclaration() bool {
	return FirstDecl <= ck && ck <= LastDecl
}

// Determine whether the given cursor kind represents an extra declaration.
func (ck Kind) IsExtraDeclaration() bool {
	return FirstExtraDecl <= ck && ck <= LastExtraDecl
}

// Determine whether the given cursor kind represents a declaration.
func (ck Kind) IsDeclaration() bool {
	return ck.IsNormDeclaration() || ck.IsExtraDeclaration()
}

/*
	Determine whether the given cursor kind represents a simple
	reference.

	Note that other kinds of cursors (such as expressions) can also refer to
	other cursors. Use clang_getCursorReferenced() to determine whether a
	particular cursor refers to another entity.
*/
func (ck Kind) IsReference() bool {
	return FirstRef <= ck && ck <= LastRef
}

// Determine whether the given cursor kind represents an expression.
func (ck Kind) IsExpression() bool {
	return FirstExpr <= ck && ck <= LastExpr
}

// Determine whether the given cursor kind represents a statement.
func (ck Kind) IsStatement() bool {
	return FirstStmt <= ck && ck <= LastStmt
}

// Determine whether the given cursor kind represents an attribute.
func (ck Kind) IsAttribute() bool {
	return FirstAttr <= ck && ck <= LastAttr
}

// Determine whether the given cursor kind represents an invalid cursor.
func (ck Kind) IsInvalid() bool {
	return FirstInvalid <= ck && ck <= LastInvalid
}

// Determine whether the given cursor kind represents a translation unit.
func (ck Kind) IsTranslationUnit() bool {
	return ck == TranslationUnit
}

// * Determine whether the given cursor represents a preprocessing element, such as a preprocessor directive or macro instantiation.
func (ck Kind) IsPreprocessing() bool {
	return FirstPreprocessing <= ck && ck <= LastPreprocessing
}

// * Determine whether the given cursor represents a currently unexposed piece of the AST (e.g., CXUnexposedStmt).
func (ck Kind) IsUnexposed() bool {
	switch ck {
	case UnexposedDecl, UnexposedExpr, UnexposedStmt, UnexposedAttr:
		return true
	}
	return false
}

// IsLiteral returns true for Literal kinds; these return no Spelling string, but getting the token is an option.
// There are some CursorKinds with the name Literal in them that sound more like expressions so there weren't
// included here.
func (ck Kind) IsLiteral() bool {
	switch ck {
	case IntegerLiteral,
		FloatingLiteral,
		ImaginaryLiteral,
		StringLiteral,
		CharacterLiteral,
		ObjCStringLiteral:
		return true
	}
	return false
}
