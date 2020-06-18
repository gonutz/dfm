package dfm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func parse(code []rune) (Object, error) {
	return newParser(code).parseObject()
}

func newParser(code []rune) *parser {
	return &parser{tokens: newTokenizer(code)}
}

type parser struct {
	tokens          tokenizer
	previewToken    token
	hasPreviewToken bool
	err             error
}

func (p *parser) parseObject() (Object, error) {
	if p.err != nil {
		return Object{}, p.err
	}

	var obj Object

	if p.peekWord("object") {
		p.word("object")
		obj.Kind = Plain
	} else if p.peekWord("inherited") {
		p.word("inherited")
		obj.Kind = Inherited
	} else if p.peekWord("inline") {
		p.word("inline")
		obj.Kind = Inline
	}

	nameOrType := p.identifier("object name (or type for anonymous objects)")
	if p.peeksAt(':') {
		obj.Name = nameOrType
		p.token(':')
		obj.Type = p.identifier("object type")
	} else {
		obj.Type = nameOrType
	}

	if p.peeksAt('[') {
		p.token('[')
		index := p.parseValue()
		if i, ok := index.(Int); ok {
			obj.HasIndex = true
			obj.Index = int(i)
		} else {
			p.err = fmt.Errorf("object index must be integer but was %#v", index)
			return obj, p.err
		}
		p.token(']')
	}

	for p.err == nil {
		if p.peekEOF() || p.peekWord("end") {
			p.nextToken()
			break
		}
		obj.Properties = append(obj.Properties, p.parseProperty())
	}
	return obj, p.err
}

func (p *parser) parseProperty() Property {
	var prop Property

	if p.peekWord("object") || p.peekWord("inherited") || p.peekWord("inline") {
		child, err := p.parseObject()
		if err != nil {
			p.err = err
			return prop
		}
		prop.Name = child.Name
		prop.Name = child.Name
		prop.Value = child
	} else {
		prop.Name = p.identifier("property name")
		for p.peekToken().tokenType == '.' {
			p.nextToken()
			prop.Name += "." + p.identifier("property name")
		}
		p.token('=')
		prop.Value = p.parseValue()
	}

	return prop
}

func (p *parser) peekEOF() bool {
	return p.peekToken().tokenType == tokenEOF
}

func (p *parser) peekWord(text string) bool {
	t := p.peekToken()
	return t.tokenType == tokenWord && strings.ToLower(t.text) == text
}

func (p *parser) parseValue() PropertyValue {
	if p.err != nil {
		return nil
	}

	switch p.peekToken().tokenType {
	case '+', '-', tokenInteger, tokenFloat:
		sign := 1
		t := p.nextToken()
		for t.tokenType == '+' || t.tokenType == '-' {
			if t.tokenType == '-' {
				sign *= -1
			}
			t = p.nextToken()
		}

		if t.tokenType == tokenInteger {
			n, err := strconv.Atoi(t.text)
			if err != nil {
				p.err = fmt.Errorf("error parsing integer literal: %v", err)
			}
			return Int(sign * n)
		} else {
			n, err := strconv.ParseFloat(t.text, 64)
			if err != nil {
				p.err = fmt.Errorf("error parsing floating point literal: %v", err)
			}
			return Float(float64(sign) * n)
		}
	case tokenCharacter, tokenString:
		// String literals can span multiple lines if connected with a '+'.
		// String literals can contain quoted strings and character literals.
		//
		//     'string'#9'with tab' +
		//     #13#10'second line#13#10
		//     'This is another string (no ''+'' before it)'
		//
		// A line break terminates a string if the last token was not a '+'.
		tokenToString := func(t token) string {
			if t.tokenType == tokenString {
				return escapeString(t.text)
			} else if t.tokenType == tokenCharacter {
				n, err := strconv.Atoi(t.text[1:])
				if err != nil {
					p.err = fmt.Errorf("error parsing character literal: %v", err)
					return ""
				}
				return string(rune(n))
			}
			return ""
		}

		// Start this string with the string/character token.
		t := p.nextToken()
		s := tokenToString(t)
		lastLine := t.line
		lastWasPlus := false
		// The string continues if other strings/characters appear on the same
		// line or if the last token on a line is a '+', in that case the next
		// string/character token continues the string.
		for {
			next := p.peekToken()
			if next.tokenType == '+' {
				lastWasPlus = true
				p.nextToken()
			} else if (next.tokenType == tokenString || next.tokenType == tokenCharacter) &&
				(lastWasPlus || next.line == lastLine) {
				t := p.nextToken()
				s += tokenToString(t)
				lastLine = t.line
				lastWasPlus = false
			} else {
				break
			}
		}
		return String(s)
	case '[':
		p.nextToken()
		var set Set
		for {
			if p.peekEOF() {
				p.err = errors.New("premature EOF in set")
				return nil
			}

			if p.peeksAt(']') {
				p.nextToken()
				break
			}
			if p.peeksAt(',') {
				p.nextToken()
			}
			set = append(set, p.parseValue())
		}
		return set
	case '(':
		p.nextToken()
		var tuple Tuple
		for {
			if p.peekEOF() {
				p.err = errors.New("premature EOF in tuple")
				return nil
			}

			if p.peeksAt(')') {
				p.nextToken()
				break
			}
			tuple = append(tuple, p.parseValue())
		}
		return tuple
	case '{':
		p.nextToken()
		code := p.tokens.findClosingBrace()
		p.token('}')

		// Copy all hex bytes to the start of code.
		hex := func(r rune) bool {
			return '0' <= r && r <= '9' ||
				'A' <= r && r <= 'F' ||
				'a' <= r && r <= 'f'
		}
		n := 0
		for i := range code {
			if hex(code[i]) {
				code[n] = code[i]
				n++
			}
		}
		code = code[:n]

		// Convert all ASCII hex pairs to bytes.
		b := make([]byte, len(code)/2)
		unHex := func(b rune) byte {
			if '0' <= b && b <= '9' {
				return byte(b) - '0'
			}
			if 'a' <= b && b <= 'f' {
				return 10 + byte(b) - 'a'
			}
			return 10 + byte(b) - 'A'
		}
		for i := range b {
			b[i] = unHex(code[i*2])<<4 | unHex(code[i*2+1])
		}

		return Bytes(b)
	case '<':
		p.nextToken()
		var items Items
		for !p.peeksAt('>') && p.err == nil {
			p.word("item")
			var item []Property
			for !p.peekWord("end") {
				item = append(item, p.parseProperty())
			}
			p.word("end")
			items = append(items, item)
		}
		p.nextToken() // Skip '>'.
		return items
	case tokenWord:
		t := p.nextToken()
		text := strings.ToLower(t.text)
		if text == "false" {
			return Bool(false)
		} else if text == "true" {
			return Bool(true)
		} else {
			id := t.text
			for p.peeksAt('.') {
				p.nextToken()
				t = p.nextToken()
				if t.tokenType != tokenWord {
					p.err = fmt.Errorf("another identifier is expected after '.' but was " + t.String())
					return nil
				}
				id += "." + t.text
			}
			return Identifier(id)
		}
	default:
		p.err = fmt.Errorf("unexpected token for property value: %v", p.nextToken())
		return nil
	}
}

func escapeString(s string) string {
	return strings.Replace(
		strings.TrimSuffix(strings.TrimPrefix(s, "'"), "'"),
		"''", "'", -1,
	)
}

func (p *parser) word(text string) {
	if p.err == nil {
		t := p.nextToken()
		if t.tokenType != tokenWord || strings.ToLower(t.text) != text {
			p.err = fmt.Errorf("%q expected but was %q", text, t.text)
		}
	}
}

func (p *parser) identifier(desc string) string {
	if p.err == nil {
		t := p.nextToken()
		if t.tokenType == tokenWord {
			return t.text
		} else {
			p.err = fmt.Errorf("identifier expected as "+desc+" but was %v", t)
		}
	}
	return ""
}

func (p *parser) token(typ tokenType) {
	if p.err == nil {
		t := p.nextToken()
		if t.tokenType != typ {
			p.err = fmt.Errorf("%v expected but was %v", typ, t)
		}
	}
}

func (p *parser) peekToken() token {
	if p.err != nil {
		return token{tokenType: tokenEOF}
	}
	if !p.hasPreviewToken {
		p.previewToken = p.tokens.next()
		for p.previewToken.tokenType == tokenWhiteSpace {
			p.previewToken = p.tokens.next()
		}
		p.hasPreviewToken = true
	}
	return p.previewToken
}

func (p *parser) peeksAt(typ tokenType) bool {
	return p.peekToken().tokenType == typ
}

func (p *parser) nextToken() token {
	if p.err != nil {
		return token{tokenType: tokenEOF}
	}
	if p.hasPreviewToken {
		p.hasPreviewToken = false
		return p.previewToken
	}
	t := p.tokens.next()
	for t.tokenType == tokenWhiteSpace {
		t = p.tokens.next()
	}
	if t.tokenType == tokenIllegal {
		p.err = fmt.Errorf("Illegal token encountered: %q", t.text)
	}
	return t
}
