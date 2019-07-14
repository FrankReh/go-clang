package tls

type Kind uint32

//go:generate stringer -type=Kind

/*
	Describe the "thread-local storage (TLS) kind" of the declaration referred
	to by a cursor.
*/
const (
	None    Kind = 0
	Dynamic Kind = 1
	Static  Kind = 2
)

func Validate(i uint32) (Kind, error) {
	switch {
	case uint32(None) <= i && i <= uint32(Static):
		return Kind(i), nil
	default:
		return 0, InvalidErr
	}
}

func MustValidate(i uint32) Kind {
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
