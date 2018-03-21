package tokenkind

//
// Describes a kind of token.
//

type Kind int

//go:generate stringer -type=Kind

const (

	/**
	 * A token that contains some kind of punctuation.
	 */
	Punctuation Kind = iota

	/**
	 * A language keyword.
	 */
	Keyword

	/**
	 * An identifier (that is not a keyword).
	 */
	Identifier

	/**
	 * A numeric, string, or character literal.
	 */
	Literal

	/**
	 * A comment.
	 */
	Comment
)

func Validate(i int) (Kind, error) {
	switch {
	case 0 <= i && i <= int(Comment):
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
