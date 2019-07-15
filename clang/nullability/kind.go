package nullability

// The nullability kind of a pointer type.
type Kind int

//go:generate stringer -type=Kind

const (
	// Values of this type can never be null.
	NonNull Kind = 0

	// Values of this type can be null.
	Nullable Kind = 1

	// Whether values of this type can be null is (explicitly)
	// unspecified. This captures a (fairly rare) case where we
	// can't conclude anything about the nullability of the type even
	// though it has been considered.
	Unspecified Kind = 2

	// Nullability is not applicable to this type.
	Invalid Kind = 3
)

func Validate(i int) (Kind, error) {
	switch {
	case 0 <= i && i <= int(Invalid):
		return Kind(i), nil
	default:
		return Invalid, InvalidErr
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
