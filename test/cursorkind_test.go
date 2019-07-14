package clang_test

import (
	"testing"

	"github.com/frankreh/go-clang/clang"
	"github.com/frankreh/go-clang/clang/cursorkind"
)

var cursorKinds = []cursorkind.Kind{

	cursorkind.UnexposedDecl,
	cursorkind.StructDecl,
	cursorkind.UnionDecl,
	cursorkind.ClassDecl,
	cursorkind.EnumDecl,
	cursorkind.FieldDecl,
	cursorkind.EnumConstantDecl,
	cursorkind.FunctionDecl,
	cursorkind.VarDecl,
	cursorkind.ParmDecl,
	cursorkind.ObjCInterfaceDecl,
	cursorkind.ObjCCategoryDecl,
	cursorkind.ObjCProtocolDecl,
	cursorkind.ObjCPropertyDecl,
	cursorkind.ObjCIvarDecl,
	cursorkind.ObjCInstanceMethodDecl,
	cursorkind.ObjCClassMethodDecl,
	cursorkind.ObjCImplementationDecl,
	cursorkind.ObjCCategoryImplDecl,
	cursorkind.TypedefDecl,
	cursorkind.CXXMethod,
	cursorkind.Namespace,
	cursorkind.LinkageSpec,
	cursorkind.Constructor,
	cursorkind.Destructor,
	cursorkind.ConversionFunction,
	cursorkind.TemplateTypeParameter,
	cursorkind.NonTypeTemplateParameter,
	cursorkind.TemplateTemplateParameter,
	cursorkind.FunctionTemplate,
	cursorkind.ClassTemplate,
	cursorkind.ClassTemplatePartialSpecialization,
	cursorkind.NamespaceAlias,
	cursorkind.UsingDirective,
	cursorkind.UsingDeclaration,
	cursorkind.TypeAliasDecl,
	cursorkind.ObjCSynthesizeDecl,
	cursorkind.ObjCDynamicDecl,
	cursorkind.CXXAccessSpecifier,
	cursorkind.ObjCSuperClassRef,
	cursorkind.ObjCProtocolRef,
	cursorkind.ObjCClassRef,
	cursorkind.TypeRef,
	cursorkind.CXXBaseSpecifier,
	cursorkind.TemplateRef,
	cursorkind.NamespaceRef,
	cursorkind.MemberRef,
	cursorkind.LabelRef,
	cursorkind.OverloadedDeclRef,
	cursorkind.VariableRef,
	cursorkind.InvalidFile,
	cursorkind.NoDeclFound,
	cursorkind.NotImplemented,
	cursorkind.InvalidCode,
	cursorkind.UnexposedExpr,
	cursorkind.DeclRefExpr,
	cursorkind.MemberRefExpr,
	cursorkind.CallExpr,
	cursorkind.ObjCMessageExpr,
	cursorkind.BlockExpr,
	cursorkind.IntegerLiteral,
	cursorkind.FloatingLiteral,
	cursorkind.ImaginaryLiteral,
	cursorkind.StringLiteral,
	cursorkind.CharacterLiteral,
	cursorkind.ParenExpr,
	cursorkind.UnaryOperator,
	cursorkind.ArraySubscriptExpr,
	cursorkind.BinaryOperator,
	cursorkind.CompoundAssignOperator,
	cursorkind.ConditionalOperator,
	cursorkind.CStyleCastExpr,
	cursorkind.CompoundLiteralExpr,
	cursorkind.InitListExpr,
	cursorkind.AddrLabelExpr,
	cursorkind.StmtExpr,
	cursorkind.GenericSelectionExpr,
	cursorkind.GNUNullExpr,
	cursorkind.CXXStaticCastExpr,
	cursorkind.CXXDynamicCastExpr,
	cursorkind.CXXReinterpretCastExpr,
	cursorkind.CXXConstCastExpr,
	cursorkind.CXXFunctionalCastExpr,
	cursorkind.CXXTypeidExpr,
	cursorkind.CXXBoolLiteralExpr,
	cursorkind.CXXNullPtrLiteralExpr,
	cursorkind.CXXThisExpr,
	cursorkind.CXXThrowExpr,
	cursorkind.CXXNewExpr,
	cursorkind.CXXDeleteExpr,
	cursorkind.UnaryExpr,
	cursorkind.ObjCStringLiteral,
	cursorkind.ObjCEncodeExpr,
	cursorkind.ObjCSelectorExpr,
	cursorkind.ObjCProtocolExpr,
	cursorkind.ObjCBridgedCastExpr,
	cursorkind.PackExpansionExpr,
	cursorkind.SizeOfPackExpr,
	cursorkind.LambdaExpr,
	cursorkind.ObjCBoolLiteralExpr,
	cursorkind.ObjCSelfExpr,
	cursorkind.OMPArraySectionExpr,
	cursorkind.ObjCAvailabilityCheckExpr,
	cursorkind.UnexposedStmt,
	cursorkind.LabelStmt,
	cursorkind.CompoundStmt,
	cursorkind.CaseStmt,
	cursorkind.DefaultStmt,
	cursorkind.IfStmt,
	cursorkind.SwitchStmt,
	cursorkind.WhileStmt,
	cursorkind.DoStmt,
	cursorkind.ForStmt,
	cursorkind.GotoStmt,
	cursorkind.IndirectGotoStmt,
	cursorkind.ContinueStmt,
	cursorkind.BreakStmt,
	cursorkind.ReturnStmt,
	cursorkind.GCCAsmStmt,
	cursorkind.AsmStmt,
	cursorkind.ObjCAtTryStmt,
	cursorkind.ObjCAtCatchStmt,
	cursorkind.ObjCAtFinallyStmt,
	cursorkind.ObjCAtThrowStmt,
	cursorkind.ObjCAtSynchronizedStmt,
	cursorkind.ObjCAutoreleasePoolStmt,
	cursorkind.ObjCForCollectionStmt,
	cursorkind.CXXCatchStmt,
	cursorkind.CXXTryStmt,
	cursorkind.CXXForRangeStmt,
	cursorkind.SEHTryStmt,
	cursorkind.SEHExceptStmt,
	cursorkind.SEHFinallyStmt,
	cursorkind.MSAsmStmt,
	cursorkind.NullStmt,
	cursorkind.DeclStmt,
	cursorkind.OMPParallelDirective,
	cursorkind.OMPSimdDirective,
	cursorkind.OMPForDirective,
	cursorkind.OMPSectionsDirective,
	cursorkind.OMPSectionDirective,
	cursorkind.OMPSingleDirective,
	cursorkind.OMPParallelForDirective,
	cursorkind.OMPParallelSectionsDirective,
	cursorkind.OMPTaskDirective,
	cursorkind.OMPMasterDirective,
	cursorkind.OMPCriticalDirective,
	cursorkind.OMPTaskyieldDirective,
	cursorkind.OMPBarrierDirective,
	cursorkind.OMPTaskwaitDirective,
	cursorkind.OMPFlushDirective,
	cursorkind.SEHLeaveStmt,
	cursorkind.OMPOrderedDirective,
	cursorkind.OMPAtomicDirective,
	cursorkind.OMPForSimdDirective,
	cursorkind.OMPParallelForSimdDirective,
	cursorkind.OMPTargetDirective,
	cursorkind.OMPTeamsDirective,
	cursorkind.OMPTaskgroupDirective,
	cursorkind.OMPCancellationPointDirective,
	cursorkind.OMPCancelDirective,
	cursorkind.OMPTargetDataDirective,
	cursorkind.OMPTaskLoopDirective,
	cursorkind.OMPTaskLoopSimdDirective,
	cursorkind.OMPDistributeDirective,
	cursorkind.OMPTargetEnterDataDirective,
	cursorkind.OMPTargetExitDataDirective,
	cursorkind.OMPTargetParallelDirective,
	cursorkind.OMPTargetParallelForDirective,
	cursorkind.OMPTargetUpdateDirective,
	cursorkind.OMPDistributeParallelForDirective,
	cursorkind.OMPDistributeParallelForSimdDirective,
	cursorkind.OMPDistributeSimdDirective,
	cursorkind.OMPTargetParallelForSimdDirective,
	cursorkind.OMPTargetSimdDirective,
	cursorkind.OMPTeamsDistributeDirective,
	cursorkind.OMPTeamsDistributeSimdDirective,
	cursorkind.OMPTeamsDistributeParallelForSimdDirective,
	cursorkind.OMPTeamsDistributeParallelForDirective,
	cursorkind.OMPTargetTeamsDirective,
	cursorkind.OMPTargetTeamsDistributeDirective,
	cursorkind.OMPTargetTeamsDistributeParallelForDirective,
	cursorkind.OMPTargetTeamsDistributeParallelForSimdDirective,
	cursorkind.OMPTargetTeamsDistributeSimdDirective,
	cursorkind.TranslationUnit,
	cursorkind.UnexposedAttr,
	cursorkind.IBActionAttr,
	cursorkind.IBOutletAttr,
	cursorkind.IBOutletCollectionAttr,
	cursorkind.CXXFinalAttr,
	cursorkind.CXXOverrideAttr,
	cursorkind.AnnotateAttr,
	cursorkind.AsmLabelAttr,
	cursorkind.PackedAttr,
	cursorkind.PureAttr,
	cursorkind.ConstAttr,
	cursorkind.NoDuplicateAttr,
	cursorkind.CUDAConstantAttr,
	cursorkind.CUDADeviceAttr,
	cursorkind.CUDAGlobalAttr,
	cursorkind.CUDAHostAttr,
	cursorkind.CUDASharedAttr,
	cursorkind.VisibilityAttr,
	cursorkind.DLLExport,
	cursorkind.DLLImport,
	cursorkind.PreprocessingDirective,
	cursorkind.MacroDefinition,
	cursorkind.MacroExpansion,
	cursorkind.MacroInstantiation,
	cursorkind.InclusionDirective,
	cursorkind.ModuleImportDecl,
	cursorkind.TypeAliasTemplateDecl,
	cursorkind.StaticAssert,
	cursorkind.FriendDecl,
	cursorkind.OverloadCandidate,
}

type compareCall struct {
	name string
	fast func(cursorkind.Kind) bool
	slow func(cursorkind.Kind) bool
}

var compareCalls = []compareCall{
	{
		name: "IsDeclaration",
		fast: cursorkind.Kind.IsDeclaration,
		slow: clang.IsDeclaration,
	},
	{
		name: "IsReference",
		fast: cursorkind.Kind.IsReference,
		slow: clang.IsReference,
	},
	{
		name: "IsExpression",
		fast: cursorkind.Kind.IsExpression,
		slow: clang.IsExpression,
	},
	{
		name: "IsStatement",
		fast: cursorkind.Kind.IsStatement,
		slow: clang.IsStatement,
	},
	{
		name: "IsAttribute",
		fast: cursorkind.Kind.IsAttribute,
		slow: clang.IsAttribute,
	},
	{
		name: "IsInvalid",
		fast: cursorkind.Kind.IsInvalid,
		slow: clang.IsInvalid,
	},
	{
		name: "IsTranslationUnit",
		fast: cursorkind.Kind.IsTranslationUnit,
		slow: clang.IsTranslationUnit,
	},
	{
		name: "IsPreprocessing",
		fast: cursorkind.Kind.IsPreprocessing,
		slow: clang.IsPreprocessing,
	},
	{
		name: "IsUnexposed",
		fast: cursorkind.Kind.IsUnexposed,
		slow: clang.IsUnexposed,
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
