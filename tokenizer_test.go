package dfm

import (
	"testing"

	"github.com/gonutz/check"
)

func TestTokenize(t *testing.T) {
	checkTokens(t,
		`object Dialog: TDialog
  Left.X=0
  +99
  -11
  [a,b]
  (123 456)
  'string'
  'quoted '' string'
  #32
  {123F}
  1.5
  1e5
  1.5E-2
  1E+1
  <>
end`,
		tok(tokenWord, "object"),
		tok(tokenWhiteSpace, " "),
		tok(tokenWord, "Dialog"),
		tok(':', ":"),
		tok(tokenWhiteSpace, " "),
		tok(tokenWord, "TDialog"),
		tok(tokenWhiteSpace, "\n  "),
		tok(tokenWord, "Left"),
		tok('.', "."),
		tok(tokenWord, "X"),
		tok('=', "="),
		tok(tokenInteger, "0"),
		tok(tokenWhiteSpace, "\n  "),
		tok('+', "+"),
		tok(tokenInteger, "99"),
		tok(tokenWhiteSpace, "\n  "),
		tok('-', "-"),
		tok(tokenInteger, "11"),
		tok(tokenWhiteSpace, "\n  "),
		tok('[', "["),
		tok(tokenWord, "a"),
		tok(',', ","),
		tok(tokenWord, "b"),
		tok(']', "]"),
		tok(tokenWhiteSpace, "\n  "),
		tok('(', "("),
		tok(tokenInteger, "123"),
		tok(tokenWhiteSpace, " "),
		tok(tokenInteger, "456"),
		tok(')', ")"),
		tok(tokenWhiteSpace, "\n  "),
		tok(tokenString, "'string'"),
		tok(tokenWhiteSpace, "\n  "),
		tok(tokenString, "'quoted '' string'"),
		tok(tokenWhiteSpace, "\n  "),
		tok(tokenCharacter, "#32"),
		tok(tokenWhiteSpace, "\n  "),
		tok('{', "{"),
		tok(tokenInteger, "123"),
		tok(tokenWord, "F"),
		tok('}', "}"),
		tok(tokenWhiteSpace, "\n  "),
		tok(tokenFloat, "1.5"),
		tok(tokenWhiteSpace, "\n  "),
		tok(tokenFloat, "1e5"),
		tok(tokenWhiteSpace, "\n  "),
		tok(tokenFloat, "1.5E-2"),
		tok(tokenWhiteSpace, "\n  "),
		tok(tokenFloat, "1E+1"),
		tok(tokenWhiteSpace, "\n  "),
		tok('<', "<"),
		tok('>', ">"),
		tok(tokenWhiteSpace, "\n"),
		tok(tokenWord, "end"),
	)
}

func TestTokenPositions(t *testing.T) {
	tokens := tokenize(`object X
  Left
end`)
	check.Eq(t, len(tokens), 7)
	check.Eq(t, tokens[0], token{
		tokenType: tokenWord,
		text:      "object",
		line:      1,
		col:       1,
	})
	check.Eq(t, tokens[1], token{
		tokenType: tokenWhiteSpace,
		text:      " ",
		line:      1,
		col:       7,
	})
	check.Eq(t, tokens[2], token{
		tokenType: tokenWord,
		text:      "X",
		line:      1,
		col:       8,
	})
	check.Eq(t, tokens[3], token{
		tokenType: tokenWhiteSpace,
		text:      "\n  ",
		line:      1,
		col:       9,
	})
	check.Eq(t, tokens[4], token{
		tokenType: tokenWord,
		text:      "Left",
		line:      2,
		col:       3,
	})
	check.Eq(t, tokens[5], token{
		tokenType: tokenWhiteSpace,
		text:      "\n",
		line:      2,
		col:       7,
	})
	check.Eq(t, tokens[6], token{
		tokenType: tokenWord,
		text:      "end",
		line:      3,
		col:       1,
	})
}

func tok(typ tokenType, text string) token {
	return token{tokenType: typ, text: text}
}

func checkTokens(t *testing.T, code string, want ...token) {
	t.Helper()
	have := tokenize(code)
	eq := len(want) == len(have)
	if eq {
		for i := range want {
			eq = eq &&
				want[i].tokenType == have[i].tokenType &&
				want[i].text == have[i].text
		}
	}
	if !eq {
		t.Error(printTokenComparison(want, have))
	}
}

func tokenize(code string) []token {
	lex := newTokenizer([]rune(code))
	var tokens []token
	for {
		if t := lex.next(); t.tokenType != tokenEOF {
			tokens = append(tokens, t)
		} else {
			break
		}
	}
	return tokens
}

func printTokenComparison(want, have []token) string {
	n := len(want)
	if len(have) > n {
		n = len(have)
	}
	s := "want <-> have"
	for i := 0; i < n; i++ {
		s += "\n"

		if i < len(want) && i < len(have) &&
			want[i].tokenType == have[i].tokenType &&
			want[i].text == have[i].text {
			s += "  "
		} else {
			s += "x "
		}

		if i < len(want) {
			s += want[i].String() + " "
		} else {
			s += "no more tokens"
		}

		s += "<->"

		if i < len(have) {
			s += " " + have[i].String()
		} else {
			s += "no more tokens"
		}
	}
	return s
}
