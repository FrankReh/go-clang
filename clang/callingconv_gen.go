package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

// Describes the calling convention of a function type
type CallingConv uint32

const (
	CallingConv_Default           CallingConv = C.CXCallingConv_Default
	CallingConv_C                 CallingConv = C.CXCallingConv_C
	CallingConv_X86StdCall        CallingConv = C.CXCallingConv_X86StdCall
	CallingConv_X86FastCall       CallingConv = C.CXCallingConv_X86FastCall
	CallingConv_X86ThisCall       CallingConv = C.CXCallingConv_X86ThisCall
	CallingConv_X86Pascal         CallingConv = C.CXCallingConv_X86Pascal
	CallingConv_AAPCS             CallingConv = C.CXCallingConv_AAPCS
	CallingConv_AAPCS_VFP         CallingConv = C.CXCallingConv_AAPCS_VFP
	CallingConv_X86RegCall        CallingConv = C.CXCallingConv_X86RegCall
	CallingConv_IntelOclBicc      CallingConv = C.CXCallingConv_IntelOclBicc
	CallingConv_Win64             CallingConv = C.CXCallingConv_Win64
	CallingConv_X86_64Win64       CallingConv = C.CXCallingConv_X86_64Win64
	CallingConv_X86_64SysV        CallingConv = C.CXCallingConv_X86_64SysV
	CallingConv_X86VectorCall     CallingConv = C.CXCallingConv_X86VectorCall
	CallingConv_Swift             CallingConv = C.CXCallingConv_Swift
	CallingConv_PreserveMost      CallingConv = C.CXCallingConv_PreserveMost
	CallingConv_PreserveAll       CallingConv = C.CXCallingConv_PreserveAll
	CallingConv_AArch64VectorCall CallingConv = C.CXCallingConv_AArch64VectorCall

	CallingConv_Invalid   CallingConv = C.CXCallingConv_Invalid
	CallingConv_Unexposed CallingConv = C.CXCallingConv_Unexposed
)
