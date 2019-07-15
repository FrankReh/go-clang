package printing

type Property uint32

//go:generate stringer -type=Property

/*
	Properties for the printing policy.

	See clang::PrintingPolicy for more information.
*/
const (
	Indentation Property = iota
	SuppressSpecifiers
	SuppressTagKeyword
	IncludeTagDefinition
	SuppressScope
	SuppressUnwrittenScope
	SuppressInitializers
	ConstantArraySizeAsWritten
	AnonymousTagLocations
	SuppressStrongLifetime
	SuppressLifetimeQualifiers
	SuppressTemplateArgsInCXXConstructors
	Bool
	Restrict
	Alignof
	UnderscoreAlignof
	UseVoidForZeroParams
	TerseOutput
	PolishForDeclaration
	Half
	MSWChar
	IncludeNewlines
	MSVCFormatting
	ConstantsAsWritten
	SuppressImplicitBase
	FullyQualifiedName
)

func Validate(i uint32) (Property, error) {
	switch {
	case i <= uint32(FullyQualifiedName):
		return Property(i), nil
	default:
		return 0, InvalidErr
	}
}

func MustValidate(i uint32) Property {
	property, err := Validate(i)
	if err != nil {
		panic(err.Error() + " " + Property(i).String())
	}
	return property
}

type Error string

const InvalidErr Error = "Invalid value"

func (e Error) Error() string {
	return string(e)
}
