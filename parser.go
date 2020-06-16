package dfm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func newParser(code string) *parser {
	return &parser{tokens: tokenizer{code: []rune(code)}}
}

type parser struct {
	tokens          tokenizer
	previewToken    token
	hasPreviewToken bool
	err             error
}

func (p *parser) parseObject() (Object, error) {
	var obj Object
	p.word("object")
	obj.Name = p.identifier("object name")
	p.token(':')
	obj.Type = p.identifier("object type")
	for {
		if p.peekEOF() || p.peekWord("end") {
			p.nextToken()
			break
		}

		var prop Property

		if p.peekWord("object") {
			child, err := p.parseObject()
			if err != nil {
				return obj, err
			}
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

		obj.Properties = append(obj.Properties, prop)
	}
	return obj, p.err
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
		var s string
		for p.peeksAt(tokenString) || p.peeksAt(tokenCharacter) || p.peeksAt('+') {
			t := p.nextToken()
			if t.tokenType == tokenString {
				s += escapeString(t.text)
			} else if t.tokenType == tokenCharacter {
				n, err := strconv.Atoi(t.text[1:])
				if err != nil {
					p.err = fmt.Errorf("error parsing character literal: %v", err)
					return nil
				}
				s += string(rune(n))
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
	case tokenWord:
		t := p.nextToken()
		text := strings.ToLower(t.text)
		if text == "false" {
			return Bool(false)
		} else if text == "true" {
			return Bool(true)
		} else {
			return Identifier(t.text)
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
