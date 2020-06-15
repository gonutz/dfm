package dfm

import "fmt"

type token struct {
	tokenType tokenType
	text      string
}

// tokenType is a rune because single characters are used directly as their
// token type, e.g. ',' '+' or ':'.
type tokenType rune

const (
	tokenIllegal    tokenType = -1
	tokenEOF        tokenType = 0
	tokenWhiteSpace tokenType = 256
	tokenWord       tokenType = 257
	tokenInteger    tokenType = 258
	tokenString     tokenType = 259
	tokenCharacter  tokenType = 260
	tokenFloat      tokenType = 261
)

func (t token) String() string {
	return fmt.Sprintf("%v: %q", t.tokenType, t.text)
}

func (t tokenType) String() string {
	switch t {
	case tokenIllegal:
		return "illegal token"
	case tokenEOF:
		return "end of file"
	case tokenWhiteSpace:
		return "white space"
	case tokenWord:
		return "word"
	case tokenInteger:
		return "integer"
	case tokenString:
		return "string"
	case tokenCharacter:
		return "character"
	case tokenFloat:
		return "floating point number"
	default:
		return fmt.Sprintf("token %q (%d)", string(t), int(t))
	}
}
