package dfm

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
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
		p.write(p.indent, o.Kind.String(), " ", o.Type)
	} else {
		p.write(p.indent, o.Kind.String(), " ", o.Name, ": ", o.Type)
	}
	if o.HasIndex {
		p.write(" [", strconv.Itoa(o.Index), "]")
	}
	p.WriteString("\r\n")
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
		// Delphi prints floating point numbers with 18 digits after the dot,
		// always. When converting to a string it uses the "best" visual style,
		// this is the -1 for strconv.FormatFloat in Go, and it just fills it
		// with 0s. To be as close to Delphi as possible, we do the same here.
		//
		// Also Delphi uses E notation for numbers above 1e+16. Below that, even
		// for numbers below 1e-16, it uses digits only.
		var s string
		if f >= 1e+16 {
			s = strconv.FormatFloat(f, 'e', -1, 64)
		} else {
			s = strconv.FormatFloat(f, 'f', -1, 64)
		}
		if strings.Contains(s, "e") {
			s = strings.Replace(s, "e+", "E", 1)
			s = strings.Replace(s, "e-", "E-", 1)
		} else {
			dot := strings.Index(s, ".")
			if dot == -1 {
				s += ".000000000000000000"
			} else {
				zeros := 18 - (len(s) - 1 - dot)
				for i := 0; i < zeros; i++ {
					s += "0"
				}
			}
		}
		p.WriteString(s)
	case Bool:
		if v {
			p.WriteString("True")
		} else {
			p.WriteString("False")
		}
	case String:
		s := string(v)
		if s == "" {
			// The empty string is a special case that is not handled by the
			// below logic.
			p.WriteString("''")
		}

		const maxLineLen = 64
		lineLen := 0
		oneLine := utf8.RuneCountInString(s) <= maxLineLen
		if !oneLine {
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

		for _, r := range s {
			if lineLen >= maxLineLen {
				beInString(false)
				p.write(" +\r\n", p.indent)
				lineLen = 0
			}

			if 32 <= r && r < 127 && r != '\'' {
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

		if !oneLine {
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
