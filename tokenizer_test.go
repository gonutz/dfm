package dfm

import "testing"

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
		tok(tokenWhiteSpace, "\n"),
		tok(tokenWord, "end"),
	)
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
	lex := tokenizer{code: []rune(code)}
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

		if i < len(want) && i < len(have) && want[i] == have[i] {
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
