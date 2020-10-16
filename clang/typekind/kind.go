package typekind

type Kind int

//go:generate stringer -type=Kind

const (
	/**
	 * Represents an invalid type (e.g., where no type is available).
	 */
	Invalid Kind = 0

	/**
	 * A type whose specific kind is not exposed via this
	 * interface.
	 */
	Unexposed Kind = 1

	/* Builtin types */
	Void         Kind = 2
	Bool         Kind = 3
	Char_U       Kind = 4
	UChar        Kind = 5
	Char16       Kind = 6
	Char32       Kind = 7
	UShort       Kind = 8
	UInt         Kind = 9
	ULong        Kind = 10
	ULongLong    Kind = 11
	UInt128      Kind = 12
	Char_S       Kind = 13
	SChar        Kind = 14
	WChar        Kind = 15
	Short        Kind = 16
	Int          Kind = 17
	Long         Kind = 18
	LongLong     Kind = 19
	Int128       Kind = 20
	Float        Kind = 21
	Double       Kind = 22
	LongDouble   Kind = 23
	NullPtr      Kind = 24
	Overload     Kind = 25
	Dependent    Kind = 26
	ObjCId       Kind = 27
	ObjCClass    Kind = 28
	ObjCSel      Kind = 29
	Float128     Kind = 30
	Half         Kind = 31
	Float16      Kind = 32
	ShortAccum   Kind = 33
	Accum        Kind = 34
	LongAccum    Kind = 35
	UShortAccum  Kind = 36
	UAccum       Kind = 37
	ULongAccum   Kind = 38
	FirstBuiltin Kind = Void
	LastBuiltin  Kind = ULongAccum

	Complex             Kind = 100
	Pointer             Kind = 101
	BlockPointer        Kind = 102
	LValueReference     Kind = 103
	RValueReference     Kind = 104
	Record              Kind = 105
	Enum                Kind = 106
	Typedef             Kind = 107
	ObjCInterface       Kind = 108
	ObjCObjectPointer   Kind = 109
	FunctionNoProto     Kind = 110
	FunctionProto       Kind = 111
	ConstantArray       Kind = 112
	Vector              Kind = 113
	IncompleteArray     Kind = 114
	VariableArray       Kind = 115
	DependentSizedArray Kind = 116
	MemberPointer       Kind = 117
	Auto                Kind = 118

	/**
	 * Represents a type that was referred to using an elaborated type keyword.
	 *
	 * E.g., struct S, or via a qualified name, e.g., N::M::type, or both.
	 */
	Elaborated Kind = 119

	/* OpenCL PipeType. */
	Pipe Kind = 120

	/* OpenCL builtin types. */
	OCLImage1dRO               Kind = 121
	OCLImage1dArrayRO          Kind = 122
	OCLImage1dBufferRO         Kind = 123
	OCLImage2dRO               Kind = 124
	OCLImage2dArrayRO          Kind = 125
	OCLImage2dDepthRO          Kind = 126
	OCLImage2dArrayDepthRO     Kind = 127
	OCLImage2dMSAARO           Kind = 128
	OCLImage2dArrayMSAARO      Kind = 129
	OCLImage2dMSAADepthRO      Kind = 130
	OCLImage2dArrayMSAADepthRO Kind = 131
	OCLImage3dRO               Kind = 132
	OCLImage1dWO               Kind = 133
	OCLImage1dArrayWO          Kind = 134
	OCLImage1dBufferWO         Kind = 135
	OCLImage2dWO               Kind = 136
	OCLImage2dArrayWO          Kind = 137
	OCLImage2dDepthWO          Kind = 138
	OCLImage2dArrayDepthWO     Kind = 139
	OCLImage2dMSAAWO           Kind = 140
	OCLImage2dArrayMSAAWO      Kind = 141
	OCLImage2dMSAADepthWO      Kind = 142
	OCLImage2dArrayMSAADepthWO Kind = 143
	OCLImage3dWO               Kind = 144
	OCLImage1dRW               Kind = 145
	OCLImage1dArrayRW          Kind = 146
	OCLImage1dBufferRW         Kind = 147
	OCLImage2dRW               Kind = 148
	OCLImage2dArrayRW          Kind = 149
	OCLImage2dDepthRW          Kind = 150
	OCLImage2dArrayDepthRW     Kind = 151
	OCLImage2dMSAARW           Kind = 152
	OCLImage2dArrayMSAARW      Kind = 153
	OCLImage2dMSAADepthRW      Kind = 154
	OCLImage2dArrayMSAADepthRW Kind = 155
	OCLImage3dRW               Kind = 156
	OCLSampler                 Kind = 157
	OCLEvent                   Kind = 158
	OCLQueue                   Kind = 159
	OCLReserveID               Kind = 160

	ObjCObject    Kind = 161
	ObjCTypeParam Kind = 162
	Attributed    Kind = 163

	OCLIntelSubgroupAVCMcePayload                  Kind = 164
	OCLIntelSubgroupAVCImePayload                  Kind = 165
	OCLIntelSubgroupAVCRefPayload                  Kind = 166
	OCLIntelSubgroupAVCSicPayload                  Kind = 167
	OCLIntelSubgroupAVCMceResult                   Kind = 168
	OCLIntelSubgroupAVCImeResult                   Kind = 169
	OCLIntelSubgroupAVCRefResult                   Kind = 170
	OCLIntelSubgroupAVCSicResult                   Kind = 171
	OCLIntelSubgroupAVCImeResultSingleRefStreamout Kind = 172
	OCLIntelSubgroupAVCImeResultDualRefStreamout   Kind = 173
	OCLIntelSubgroupAVCImeSingleRefStreamin        Kind = 174
	OCLIntelSubgroupAVCImeDualRefStreamin          Kind = 175
	ExtVector                                      Kind = 176
)

func Validate(i int) (Kind, error) {
	switch {
	case 0 <= i && i <= int(LastBuiltin):
		return Kind(i), nil
	case int(Complex) <= i && i <= int(ExtVector):
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

func (k Kind) IsBuiltin() bool {
	return FirstBuiltin <= k && k <= LastBuiltin
}
