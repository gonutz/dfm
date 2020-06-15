package dfm

import (
	"fmt"
	"strconv"
)

// String returns the text representation of the Object as DFM code.
func (o Object) String() string {
	p := printer{}
	p.object(o)
	return p.buf
}

type printer struct {
	buf    string
	indent string
}

func (p *printer) object(o Object) {
	p.buf += p.indent + "object " + o.Name + ": " + o.Type + "\r\n"
	p.indent += "  "
	for _, prop := range o.Properties {
		if obj, ok := prop.Value.(Object); ok {
			p.object(obj)
		} else {
			p.property(prop)
		}
	}
	p.indent = p.indent[2:]
	p.buf += p.indent + "end\r\n"
}

func (p *printer) property(prop Property) {
	p.buf += p.indent + prop.Name + " = "
	p.propertyValue(prop.Value)
	p.buf += "\r\n"
}

func (p *printer) propertyValue(value PropertyValue) {
	switch v := value.(type) {
	case Int:
		p.buf += strconv.Itoa(int(v))
	case Bool:
		if v {
			p.buf += "True"
		} else {
			p.buf += "False"
		}
	case String:
		const maxLineLen = 63
		lineLen := 0
		if len(string(v)) > maxLineLen {
			p.indent += "  "
			p.buf += "\r\n" + p.indent
		}
		inString := false
		beInString := func(in bool) {
			if inString != in {
				p.buf += "'"
				inString = !inString
			}
		}

		for _, r := range string(v) {
			if lineLen > maxLineLen {
				beInString(false)
				p.buf += " +\r\n" + p.indent
				lineLen = 0
			}

			if r == '\'' {
				beInString(true)
				p.buf += "''" // Escape ' quotes.
			} else if 32 <= r && r < 127 {
				beInString(true)
				p.buf += string(r) // Printable ASCII character.
			} else {
				beInString(false)
				p.buf += "#" + strconv.Itoa(int(r)) // Unicode character.
			}
			lineLen++
		}
		beInString(false)

		if len(string(v)) > maxLineLen {
			p.indent = p.indent[2:]
		}
	case Identifier:
		p.buf += string(v)
	case Set:
		p.buf += "["
		for i := range v {
			if i > 0 {
				p.buf += ", "
			}
			p.propertyValue(v[i])
		}
		p.buf += "]"
	case Tuple:
		p.buf += "("
		p.indent += "  "
		for i := range v {
			p.buf += "\r\n" + p.indent
			p.propertyValue(v[i])
		}
		p.buf += ")"
		p.indent = p.indent[2:]
	case Bytes:
		p.indent += "  "
		p.buf += "{\r\n" + p.indent

		hexNibble := []string{
			"0", "1", "2", "3",
			"4", "5", "6", "7",
			"8", "9", "A", "B",
			"C", "D", "E", "F",
		}

		const maxLineLen = 31
		lineLen := 0
		for _, b := range v {
			if lineLen > maxLineLen {
				p.buf += "\r\n" + p.indent
				lineLen = 0
			}
			p.buf += hexNibble[b&0xF0>>4]
			p.buf += hexNibble[b&0x0F]
			lineLen++
		}

		p.buf += "}"
		p.indent = p.indent[2:]
	default:
		panic("unhandled property type " + fmt.Sprintf("%T", v))
	}
}
