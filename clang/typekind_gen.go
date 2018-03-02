package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

// Describes the kind of type
type TypeKind uint32

const (
	// Represents an invalid type (e.g., where no type is available).
	Type_Invalid TypeKind = C.CXType_Invalid
	// A type whose specific kind is not exposed via this interface.
	Type_Unexposed TypeKind = C.CXType_Unexposed
	// A type whose specific kind is not exposed via this interface.
	Type_Void TypeKind = C.CXType_Void
	// A type whose specific kind is not exposed via this interface.
	Type_Bool TypeKind = C.CXType_Bool
	// A type whose specific kind is not exposed via this interface.
	Type_Char_U TypeKind = C.CXType_Char_U
	// A type whose specific kind is not exposed via this interface.
	Type_UChar TypeKind = C.CXType_UChar
	// A type whose specific kind is not exposed via this interface.
	Type_Char16 TypeKind = C.CXType_Char16
	// A type whose specific kind is not exposed via this interface.
	Type_Char32 TypeKind = C.CXType_Char32
	// A type whose specific kind is not exposed via this interface.
	Type_UShort TypeKind = C.CXType_UShort
	// A type whose specific kind is not exposed via this interface.
	Type_UInt TypeKind = C.CXType_UInt
	// A type whose specific kind is not exposed via this interface.
	Type_ULong TypeKind = C.CXType_ULong
	// A type whose specific kind is not exposed via this interface.
	Type_ULongLong TypeKind = C.CXType_ULongLong
	// A type whose specific kind is not exposed via this interface.
	Type_UInt128 TypeKind = C.CXType_UInt128
	// A type whose specific kind is not exposed via this interface.
	Type_Char_S TypeKind = C.CXType_Char_S
	// A type whose specific kind is not exposed via this interface.
	Type_SChar TypeKind = C.CXType_SChar
	// A type whose specific kind is not exposed via this interface.
	Type_WChar TypeKind = C.CXType_WChar
	// A type whose specific kind is not exposed via this interface.
	Type_Short TypeKind = C.CXType_Short
	// A type whose specific kind is not exposed via this interface.
	Type_Int TypeKind = C.CXType_Int
	// A type whose specific kind is not exposed via this interface.
	Type_Long TypeKind = C.CXType_Long
	// A type whose specific kind is not exposed via this interface.
	Type_LongLong TypeKind = C.CXType_LongLong
	// A type whose specific kind is not exposed via this interface.
	Type_Int128 TypeKind = C.CXType_Int128
	// A type whose specific kind is not exposed via this interface.
	Type_Float TypeKind = C.CXType_Float
	// A type whose specific kind is not exposed via this interface.
	Type_Double TypeKind = C.CXType_Double
	// A type whose specific kind is not exposed via this interface.
	Type_LongDouble TypeKind = C.CXType_LongDouble
	// A type whose specific kind is not exposed via this interface.
	Type_NullPtr TypeKind = C.CXType_NullPtr
	// A type whose specific kind is not exposed via this interface.
	Type_Overload TypeKind = C.CXType_Overload
	// A type whose specific kind is not exposed via this interface.
	Type_Dependent TypeKind = C.CXType_Dependent
	// A type whose specific kind is not exposed via this interface.
	Type_ObjCId TypeKind = C.CXType_ObjCId
	// A type whose specific kind is not exposed via this interface.
	Type_ObjCClass TypeKind = C.CXType_ObjCClass
	// A type whose specific kind is not exposed via this interface.
	Type_ObjCSel TypeKind = C.CXType_ObjCSel
	// A type whose specific kind is not exposed via this interface.
	Type_Float128 TypeKind = C.CXType_Float128
	// A type whose specific kind is not exposed via this interface.
	Type_Half TypeKind = C.CXType_Half
	// A type whose specific kind is not exposed via this interface.
	Type_FirstBuiltin TypeKind = C.CXType_FirstBuiltin
	// A type whose specific kind is not exposed via this interface.
	Type_LastBuiltin TypeKind = C.CXType_LastBuiltin
	// A type whose specific kind is not exposed via this interface.
	Type_Complex TypeKind = C.CXType_Complex
	// A type whose specific kind is not exposed via this interface.
	Type_Pointer TypeKind = C.CXType_Pointer
	// A type whose specific kind is not exposed via this interface.
	Type_BlockPointer TypeKind = C.CXType_BlockPointer
	// A type whose specific kind is not exposed via this interface.
	Type_LValueReference TypeKind = C.CXType_LValueReference
	// A type whose specific kind is not exposed via this interface.
	Type_RValueReference TypeKind = C.CXType_RValueReference
	// A type whose specific kind is not exposed via this interface.
	Type_Record TypeKind = C.CXType_Record
	// A type whose specific kind is not exposed via this interface.
	Type_Enum TypeKind = C.CXType_Enum
	// A type whose specific kind is not exposed via this interface.
	Type_Typedef TypeKind = C.CXType_Typedef
	// A type whose specific kind is not exposed via this interface.
	Type_ObjCInterface TypeKind = C.CXType_ObjCInterface
	// A type whose specific kind is not exposed via this interface.
	Type_ObjCObjectPointer TypeKind = C.CXType_ObjCObjectPointer
	// A type whose specific kind is not exposed via this interface.
	Type_FunctionNoProto TypeKind = C.CXType_FunctionNoProto
	// A type whose specific kind is not exposed via this interface.
	Type_FunctionProto TypeKind = C.CXType_FunctionProto
	// A type whose specific kind is not exposed via this interface.
	Type_ConstantArray TypeKind = C.CXType_ConstantArray
	// A type whose specific kind is not exposed via this interface.
	Type_Vector TypeKind = C.CXType_Vector
	// A type whose specific kind is not exposed via this interface.
	Type_IncompleteArray TypeKind = C.CXType_IncompleteArray
	// A type whose specific kind is not exposed via this interface.
	Type_VariableArray TypeKind = C.CXType_VariableArray
	// A type whose specific kind is not exposed via this interface.
	Type_DependentSizedArray TypeKind = C.CXType_DependentSizedArray
	// A type whose specific kind is not exposed via this interface.
	Type_MemberPointer TypeKind = C.CXType_MemberPointer
	// A type whose specific kind is not exposed via this interface.
	Type_Auto TypeKind = C.CXType_Auto
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_Elaborated TypeKind = C.CXType_Elaborated
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_Pipe TypeKind = C.CXType_Pipe
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage1dRO TypeKind = C.CXType_OCLImage1dRO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage1dArrayRO TypeKind = C.CXType_OCLImage1dArrayRO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage1dBufferRO TypeKind = C.CXType_OCLImage1dBufferRO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dRO TypeKind = C.CXType_OCLImage2dRO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dArrayRO TypeKind = C.CXType_OCLImage2dArrayRO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dDepthRO TypeKind = C.CXType_OCLImage2dDepthRO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dArrayDepthRO TypeKind = C.CXType_OCLImage2dArrayDepthRO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dMSAARO TypeKind = C.CXType_OCLImage2dMSAARO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dArrayMSAARO TypeKind = C.CXType_OCLImage2dArrayMSAARO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dMSAADepthRO TypeKind = C.CXType_OCLImage2dMSAADepthRO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dArrayMSAADepthRO TypeKind = C.CXType_OCLImage2dArrayMSAADepthRO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage3dRO TypeKind = C.CXType_OCLImage3dRO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage1dWO TypeKind = C.CXType_OCLImage1dWO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage1dArrayWO TypeKind = C.CXType_OCLImage1dArrayWO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage1dBufferWO TypeKind = C.CXType_OCLImage1dBufferWO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dWO TypeKind = C.CXType_OCLImage2dWO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dArrayWO TypeKind = C.CXType_OCLImage2dArrayWO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dDepthWO TypeKind = C.CXType_OCLImage2dDepthWO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dArrayDepthWO TypeKind = C.CXType_OCLImage2dArrayDepthWO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dMSAAWO TypeKind = C.CXType_OCLImage2dMSAAWO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dArrayMSAAWO TypeKind = C.CXType_OCLImage2dArrayMSAAWO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dMSAADepthWO TypeKind = C.CXType_OCLImage2dMSAADepthWO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dArrayMSAADepthWO TypeKind = C.CXType_OCLImage2dArrayMSAADepthWO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage3dWO TypeKind = C.CXType_OCLImage3dWO
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage1dRW TypeKind = C.CXType_OCLImage1dRW
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage1dArrayRW TypeKind = C.CXType_OCLImage1dArrayRW
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage1dBufferRW TypeKind = C.CXType_OCLImage1dBufferRW
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dRW TypeKind = C.CXType_OCLImage2dRW
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dArrayRW TypeKind = C.CXType_OCLImage2dArrayRW
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dDepthRW TypeKind = C.CXType_OCLImage2dDepthRW
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dArrayDepthRW TypeKind = C.CXType_OCLImage2dArrayDepthRW
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dMSAARW TypeKind = C.CXType_OCLImage2dMSAARW
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dArrayMSAARW TypeKind = C.CXType_OCLImage2dArrayMSAARW
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dMSAADepthRW TypeKind = C.CXType_OCLImage2dMSAADepthRW
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage2dArrayMSAADepthRW TypeKind = C.CXType_OCLImage2dArrayMSAADepthRW
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLImage3dRW TypeKind = C.CXType_OCLImage3dRW
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLSampler TypeKind = C.CXType_OCLSampler
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLEvent TypeKind = C.CXType_OCLEvent
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLQueue TypeKind = C.CXType_OCLQueue
	/*
		Represents a type that was referred to using an elaborated type keyword.

		E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	*/
	Type_OCLReserveID TypeKind = C.CXType_OCLReserveID
)

func (tk TypeKind) Spelling() string {
	return cx2GoString(C.clang_getTypeKindSpelling(C.enum_CXTypeKind(tk)))
}
