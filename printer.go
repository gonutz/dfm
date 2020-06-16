package dfm

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
)

// String returns the text representation of the Object as DFM code.
// Float values NaN and +-Infinity are printed as 0 since they are invalid in
// DFM files.
func (o Object) String() string {
	p := printer{}
	p.object(o)
	return p.String()
}

type printer struct {
	bytes.Buffer
	indent string
}

func (p *printer) write(s ...string) {
	for _, s := range s {
		p.WriteString(s)
	}
}

func (p *printer) incIndent() {
	p.indent += "  "
}

func (p *printer) decIndent() {
	p.indent = p.indent[:len(p.indent)-2]
}

func (p *printer) object(o Object) {
	if o.Name == "" {
		// Anonymous object.
		p.write(p.indent, o.Kind.String(), " ", o.Type, "\r\n")
	} else {
		p.write(p.indent, o.Kind.String(), " ", o.Name, ": ", o.Type, "\r\n")
	}
	p.incIndent()
	for _, prop := range o.Properties {
		if obj, ok := prop.Value.(Object); ok {
			p.object(obj)
		} else {
			p.property(prop)
		}
	}
	p.decIndent()
	p.write(p.indent, "end\r\n")
}

func (p *printer) property(prop Property) {
	p.write(p.indent, prop.Name, " = ")
	p.propertyValue(prop.Value)
	p.WriteString("\r\n")
}

func (p *printer) propertyValue(value PropertyValue) {
	switch v := value.(type) {
	case Int:
		p.WriteString(strconv.Itoa(int(v)))
	case Float:
		f := float64(v)
		if math.IsNaN(f) || math.IsInf(f, 0) {
			f = 0
		}
		p.WriteString(strconv.FormatFloat(f, 'f', 18, 64))
	case Bool:
		if v {
			p.WriteString("True")
		} else {
			p.WriteString("False")
		}
	case String:
		if string(v) == "" {
			// The empty string is a special case that is not handled by the
			// below logic.
			p.WriteString("''")
		}

		const maxLineLen = 63
		lineLen := 0
		if len(string(v)) > maxLineLen {
			p.incIndent()
			p.write("\r\n", p.indent)
		}
		inString := false
		beInString := func(in bool) {
			if inString != in {
				p.WriteByte('\'')
				inString = !inString
			}
		}

		for _, r := range string(v) {
			if lineLen > maxLineLen {
				beInString(false)
				p.write(" +\r\n", p.indent)
				lineLen = 0
			}

			if r == '\'' {
				beInString(true)
				p.WriteString("''") // Escape ' quotes.
			} else if 32 <= r && r < 127 {
				beInString(true)
				p.WriteRune(r) // Printable ASCII character.
			} else {
				beInString(false)
				// Unicode character as #<number>.
				p.WriteByte('#')
				p.WriteString(strconv.Itoa(int(r)))
			}
			lineLen++
		}
		beInString(false)

		if len(string(v)) > maxLineLen {
			p.decIndent()
		}
	case Identifier:
		p.WriteString(string(v))
	case Set:
		p.WriteByte('[')
		for i := range v {
			if i > 0 {
				p.WriteString(", ")
			}
			p.propertyValue(v[i])
		}
		p.WriteByte(']')
	case Tuple:
		p.WriteByte('(')
		p.incIndent()
		for i := range v {
			p.write("\r\n", p.indent)
			p.propertyValue(v[i])
		}
		p.WriteByte(')')
		p.decIndent()
	case Bytes:
		p.incIndent()
		p.write("{\r\n", p.indent)
		hexNibble := []byte("0123456789ABCDEF")
		const maxLineLen = 31
		lineLen := 0
		for _, b := range v {
			if lineLen > maxLineLen {
				p.write("\r\n", p.indent)
				lineLen = 0
			}
			p.WriteByte(hexNibble[b&0xF0>>4])
			p.WriteByte(hexNibble[b&0x0F])
			lineLen++
		}
		p.WriteByte('}')
		p.decIndent()
	case Items:
		p.WriteByte('<')
		p.incIndent()
		for _, properties := range v {
			p.write("\r\n", p.indent, "item\r\n")
			p.incIndent()
			for _, prop := range properties {
				p.property(prop)
			}
			p.decIndent()
			p.write(p.indent, "end")
		}
		p.decIndent()
		p.WriteByte('>')
	default:
		panic("unhandled property type " + fmt.Sprintf("%T", v))
	}
}
