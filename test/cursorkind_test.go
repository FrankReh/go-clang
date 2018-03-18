package clang_test

import (
	"testing"

	"github.com/frankreh/go-clang-v5.0/clang"
)

var cursorKinds = []clang.CursorKind{

	clang.Cursor_UnexposedDecl,
	clang.Cursor_StructDecl,
	clang.Cursor_UnionDecl,
	clang.Cursor_ClassDecl,
	clang.Cursor_EnumDecl,
	clang.Cursor_FieldDecl,
	clang.Cursor_EnumConstantDecl,
	clang.Cursor_FunctionDecl,
	clang.Cursor_VarDecl,
	clang.Cursor_ParmDecl,
	clang.Cursor_ObjCInterfaceDecl,
	clang.Cursor_ObjCCategoryDecl,
	clang.Cursor_ObjCProtocolDecl,
	clang.Cursor_ObjCPropertyDecl,
	clang.Cursor_ObjCIvarDecl,
	clang.Cursor_ObjCInstanceMethodDecl,
	clang.Cursor_ObjCClassMethodDecl,
	clang.Cursor_ObjCImplementationDecl,
	clang.Cursor_ObjCCategoryImplDecl,
	clang.Cursor_TypedefDecl,
	clang.Cursor_CXXMethod,
	clang.Cursor_Namespace,
	clang.Cursor_LinkageSpec,
	clang.Cursor_Constructor,
	clang.Cursor_Destructor,
	clang.Cursor_ConversionFunction,
	clang.Cursor_TemplateTypeParameter,
	clang.Cursor_NonTypeTemplateParameter,
	clang.Cursor_TemplateTemplateParameter,
	clang.Cursor_FunctionTemplate,
	clang.Cursor_ClassTemplate,
	clang.Cursor_ClassTemplatePartialSpecialization,
	clang.Cursor_NamespaceAlias,
	clang.Cursor_UsingDirective,
	clang.Cursor_UsingDeclaration,
	clang.Cursor_TypeAliasDecl,
	clang.Cursor_ObjCSynthesizeDecl,
	clang.Cursor_ObjCDynamicDecl,
	clang.Cursor_CXXAccessSpecifier,
	clang.Cursor_ObjCSuperClassRef,
	clang.Cursor_ObjCProtocolRef,
	clang.Cursor_ObjCClassRef,
	clang.Cursor_TypeRef,
	clang.Cursor_CXXBaseSpecifier,
	clang.Cursor_TemplateRef,
	clang.Cursor_NamespaceRef,
	clang.Cursor_MemberRef,
	clang.Cursor_LabelRef,
	clang.Cursor_OverloadedDeclRef,
	clang.Cursor_VariableRef,
	clang.Cursor_InvalidFile,
	clang.Cursor_NoDeclFound,
	clang.Cursor_NotImplemented,
	clang.Cursor_InvalidCode,
	clang.Cursor_UnexposedExpr,
	clang.Cursor_DeclRefExpr,
	clang.Cursor_MemberRefExpr,
	clang.Cursor_CallExpr,
	clang.Cursor_ObjCMessageExpr,
	clang.Cursor_BlockExpr,
	clang.Cursor_IntegerLiteral,
	clang.Cursor_FloatingLiteral,
	clang.Cursor_ImaginaryLiteral,
	clang.Cursor_StringLiteral,
	clang.Cursor_CharacterLiteral,
	clang.Cursor_ParenExpr,
	clang.Cursor_UnaryOperator,
	clang.Cursor_ArraySubscriptExpr,
	clang.Cursor_BinaryOperator,
	clang.Cursor_CompoundAssignOperator,
	clang.Cursor_ConditionalOperator,
	clang.Cursor_CStyleCastExpr,
	clang.Cursor_CompoundLiteralExpr,
	clang.Cursor_InitListExpr,
	clang.Cursor_AddrLabelExpr,
	clang.Cursor_StmtExpr,
	clang.Cursor_GenericSelectionExpr,
	clang.Cursor_GNUNullExpr,
	clang.Cursor_CXXStaticCastExpr,
	clang.Cursor_CXXDynamicCastExpr,
	clang.Cursor_CXXReinterpretCastExpr,
	clang.Cursor_CXXConstCastExpr,
	clang.Cursor_CXXFunctionalCastExpr,
	clang.Cursor_CXXTypeidExpr,
	clang.Cursor_CXXBoolLiteralExpr,
	clang.Cursor_CXXNullPtrLiteralExpr,
	clang.Cursor_CXXThisExpr,
	clang.Cursor_CXXThrowExpr,
	clang.Cursor_CXXNewExpr,
	clang.Cursor_CXXDeleteExpr,
	clang.Cursor_UnaryExpr,
	clang.Cursor_ObjCStringLiteral,
	clang.Cursor_ObjCEncodeExpr,
	clang.Cursor_ObjCSelectorExpr,
	clang.Cursor_ObjCProtocolExpr,
	clang.Cursor_ObjCBridgedCastExpr,
	clang.Cursor_PackExpansionExpr,
	clang.Cursor_SizeOfPackExpr,
	clang.Cursor_LambdaExpr,
	clang.Cursor_ObjCBoolLiteralExpr,
	clang.Cursor_ObjCSelfExpr,
	clang.Cursor_OMPArraySectionExpr,
	clang.Cursor_ObjCAvailabilityCheckExpr,
	clang.Cursor_UnexposedStmt,
	clang.Cursor_LabelStmt,
	clang.Cursor_CompoundStmt,
	clang.Cursor_CaseStmt,
	clang.Cursor_DefaultStmt,
	clang.Cursor_IfStmt,
	clang.Cursor_SwitchStmt,
	clang.Cursor_WhileStmt,
	clang.Cursor_DoStmt,
	clang.Cursor_ForStmt,
	clang.Cursor_GotoStmt,
	clang.Cursor_IndirectGotoStmt,
	clang.Cursor_ContinueStmt,
	clang.Cursor_BreakStmt,
	clang.Cursor_ReturnStmt,
	clang.Cursor_GCCAsmStmt,
	clang.Cursor_AsmStmt,
	clang.Cursor_ObjCAtTryStmt,
	clang.Cursor_ObjCAtCatchStmt,
	clang.Cursor_ObjCAtFinallyStmt,
	clang.Cursor_ObjCAtThrowStmt,
	clang.Cursor_ObjCAtSynchronizedStmt,
	clang.Cursor_ObjCAutoreleasePoolStmt,
	clang.Cursor_ObjCForCollectionStmt,
	clang.Cursor_CXXCatchStmt,
	clang.Cursor_CXXTryStmt,
	clang.Cursor_CXXForRangeStmt,
	clang.Cursor_SEHTryStmt,
	clang.Cursor_SEHExceptStmt,
	clang.Cursor_SEHFinallyStmt,
	clang.Cursor_MSAsmStmt,
	clang.Cursor_NullStmt,
	clang.Cursor_DeclStmt,
	clang.Cursor_OMPParallelDirective,
	clang.Cursor_OMPSimdDirective,
	clang.Cursor_OMPForDirective,
	clang.Cursor_OMPSectionsDirective,
	clang.Cursor_OMPSectionDirective,
	clang.Cursor_OMPSingleDirective,
	clang.Cursor_OMPParallelForDirective,
	clang.Cursor_OMPParallelSectionsDirective,
	clang.Cursor_OMPTaskDirective,
	clang.Cursor_OMPMasterDirective,
	clang.Cursor_OMPCriticalDirective,
	clang.Cursor_OMPTaskyieldDirective,
	clang.Cursor_OMPBarrierDirective,
	clang.Cursor_OMPTaskwaitDirective,
	clang.Cursor_OMPFlushDirective,
	clang.Cursor_SEHLeaveStmt,
	clang.Cursor_OMPOrderedDirective,
	clang.Cursor_OMPAtomicDirective,
	clang.Cursor_OMPForSimdDirective,
	clang.Cursor_OMPParallelForSimdDirective,
	clang.Cursor_OMPTargetDirective,
	clang.Cursor_OMPTeamsDirective,
	clang.Cursor_OMPTaskgroupDirective,
	clang.Cursor_OMPCancellationPointDirective,
	clang.Cursor_OMPCancelDirective,
	clang.Cursor_OMPTargetDataDirective,
	clang.Cursor_OMPTaskLoopDirective,
	clang.Cursor_OMPTaskLoopSimdDirective,
	clang.Cursor_OMPDistributeDirective,
	clang.Cursor_OMPTargetEnterDataDirective,
	clang.Cursor_OMPTargetExitDataDirective,
	clang.Cursor_OMPTargetParallelDirective,
	clang.Cursor_OMPTargetParallelForDirective,
	clang.Cursor_OMPTargetUpdateDirective,
	clang.Cursor_OMPDistributeParallelForDirective,
	clang.Cursor_OMPDistributeParallelForSimdDirective,
	clang.Cursor_OMPDistributeSimdDirective,
	clang.Cursor_OMPTargetParallelForSimdDirective,
	clang.Cursor_OMPTargetSimdDirective,
	clang.Cursor_OMPTeamsDistributeDirective,
	clang.Cursor_OMPTeamsDistributeSimdDirective,
	clang.Cursor_OMPTeamsDistributeParallelForSimdDirective,
	clang.Cursor_OMPTeamsDistributeParallelForDirective,
	clang.Cursor_OMPTargetTeamsDirective,
	clang.Cursor_OMPTargetTeamsDistributeDirective,
	clang.Cursor_OMPTargetTeamsDistributeParallelForDirective,
	clang.Cursor_OMPTargetTeamsDistributeParallelForSimdDirective,
	clang.Cursor_OMPTargetTeamsDistributeSimdDirective,
	clang.Cursor_TranslationUnit,
	clang.Cursor_UnexposedAttr,
	clang.Cursor_IBActionAttr,
	clang.Cursor_IBOutletAttr,
	clang.Cursor_IBOutletCollectionAttr,
	clang.Cursor_CXXFinalAttr,
	clang.Cursor_CXXOverrideAttr,
	clang.Cursor_AnnotateAttr,
	clang.Cursor_AsmLabelAttr,
	clang.Cursor_PackedAttr,
	clang.Cursor_PureAttr,
	clang.Cursor_ConstAttr,
	clang.Cursor_NoDuplicateAttr,
	clang.Cursor_CUDAConstantAttr,
	clang.Cursor_CUDADeviceAttr,
	clang.Cursor_CUDAGlobalAttr,
	clang.Cursor_CUDAHostAttr,
	clang.Cursor_CUDASharedAttr,
	clang.Cursor_VisibilityAttr,
	clang.Cursor_DLLExport,
	clang.Cursor_DLLImport,
	clang.Cursor_PreprocessingDirective,
	clang.Cursor_MacroDefinition,
	clang.Cursor_MacroExpansion,
	clang.Cursor_MacroInstantiation,
	clang.Cursor_InclusionDirective,
	clang.Cursor_ModuleImportDecl,
	clang.Cursor_TypeAliasTemplateDecl,
	clang.Cursor_StaticAssert,
	clang.Cursor_FriendDecl,
	clang.Cursor_OverloadCandidate,
}

type compareCall struct {
	name string
	fast func(clang.CursorKind) bool
	slow func(clang.CursorKind) bool
}

var compareCalls = []compareCall{
	{
		name: "IsDeclaration",
		fast: clang.CursorKind.IsDeclaration,
		slow: clang.CursorKind.IsDeclarationSlow,
	},
	{
		name: "IsReference",
		fast: clang.CursorKind.IsReference,
		slow: clang.CursorKind.IsReferenceSlow,
	},
	{
		name: "IsExpression",
		fast: clang.CursorKind.IsExpression,
		slow: clang.CursorKind.IsExpressionSlow,
	},
	{
		name: "IsStatement",
		fast: clang.CursorKind.IsStatement,
		slow: clang.CursorKind.IsStatementSlow,
	},
	{
		name: "IsAttribute",
		fast: clang.CursorKind.IsAttribute,
		slow: clang.CursorKind.IsAttributeSlow,
	},
	{
		name: "IsInvalid",
		fast: clang.CursorKind.IsInvalid,
		slow: clang.CursorKind.IsInvalidSlow,
	},
	{
		name: "IsTranslationUnit",
		fast: clang.CursorKind.IsTranslationUnit,
		slow: clang.CursorKind.IsTranslationUnitSlow,
	},
	{
		name: "IsPreprocessing",
		fast: clang.CursorKind.IsPreprocessing,
		slow: clang.CursorKind.IsPreprocessingSlow,
	},
	{
		name: "IsUnexposed",
		fast: clang.CursorKind.IsUnexposed,
		slow: clang.CursorKind.IsUnexposedSlow,
	},
}

//  TestCursorKindFastCalls tests that the fast and slow versions of
// the CursorKind boolean calls return the same results. Using the
// fast versions can save on calls into cgo.
func TestCursorKindFastCalls(t *testing.T) {
	for _, test := range compareCalls {
		t.Run(test.name, func(t *testing.T) {
			for _, cursorKind := range cursorKinds {
				fast := test.fast(cursorKind)
				slow := test.slow(cursorKind)
				if fast != slow {
					t.Errorf("%s boolean call %s, fast:%v slow:%v",
						cursorKind.String(), test.name, fast, slow)
				}
			}
		})
	}
}

func BenchmarkCursorKindFastCalls(b *testing.B) {
	for _, test := range compareCalls {
		for _, cursorKind := range cursorKinds {
			b.Run(test.name+"/fast/"+cursorKind.String(), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = test.fast(cursorKind)
				}
			})
			b.Run(test.name+"/slow/"+cursorKind.String(), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = test.slow(cursorKind)
				}
			})
			// If -short is specified, just benchmark for the first CursorKind,
			// the rest are the same anyway from my experience.
			if testing.Short() {
				break
			}
		}
	}
}
