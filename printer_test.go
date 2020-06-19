package dfm_test

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/gonutz/check"
	"github.com/gonutz/dfm"
)

var utf8bom = []byte{0xEF, 0xBB, 0xBF}

func prop(name string, value dfm.PropertyValue) dfm.Object {
	return dfm.Object{Properties: []dfm.Property{{Name: name, Value: value}}}
}

func TestASCIIobjectsAreSavedWithoutBOM(t *testing.T) {
	objs := []dfm.Object{
		{Name: "a"},
		{Type: "a"},
		prop("_", &dfm.Object{Name: "a"}),
		prop("a", dfm.Int(0)),
		prop("a", dfm.String("")),
		prop("_", dfm.Identifier("a")),
		prop("_", dfm.Set{dfm.Identifier("a")}),
		prop("_", dfm.Tuple{dfm.Identifier("a")}),
		prop("_", dfm.Items{[]dfm.Property{{Name: "a", Value: dfm.Int(0)}}}),
		prop("_", dfm.Items{[]dfm.Property{{Value: dfm.Identifier("a")}}}),

		// As long as the name is ASCII, the ä in the string literal does not
		// make a difference, it will be escaped as a character literal.
		prop("a", dfm.String("ä")),
	}
	for i, obj := range objs {
		data := obj.Print()
		check.Neq(t, data[:3], utf8bom, "object ", i)
	}
}

func TestNonASCIIobjectsAreSavedWithBOM(t *testing.T) {
	objs := []dfm.Object{
		{Name: "ä"},
		{Type: "ä"},
		prop("_", &dfm.Object{Name: "ä"}),
		prop("ä", dfm.Int(0)),
		prop("ä", dfm.String("")),
		prop("_", dfm.Identifier("ä")),
		prop("_", dfm.Set{dfm.Identifier("ä")}),
		prop("_", dfm.Tuple{dfm.Identifier("ä")}),
		prop("_", dfm.Items{[]dfm.Property{{Name: "ä", Value: dfm.Int(0)}}}),
		prop("_", dfm.Items{[]dfm.Property{{Value: dfm.Identifier("ä")}}}),
	}
	for i, obj := range objs {
		data := obj.Print()
		check.Eq(t, data[:3], utf8bom, "object ", i)
	}
}

func TestPrintDFM(t *testing.T) {
	// Things to test:
	// - if Object.Name is empty, only print the type.
	// - All Kinds of Object.
	// - Objects with index.
	// - All implementations of PropertyValue.
	// - Print Float with 18 digits after the point. The actual precision is
	//   usually lower, though (strange Delphi behavior, see printer.go).
	// - Strings with quotes and non-ASCII characters.
	// - Long Strings should span multiple lines.
	// - Identifiers with dots, both as property names and values.
	// - Floats NAN and +-INF are invalid in DFMs, default to 0.
	obj := dfm.Object{
		Name: "Dialog",
		Type: "TDialog",
		Properties: []dfm.Property{
			{Value: &dfm.Object{
				Type: "TSubObject",
				Kind: dfm.Inline,
			}},
			{Value: &dfm.Object{
				Name: "Child",
				Type: "TChild",
				Kind: dfm.Inherited,
			}},
			{Value: &dfm.Object{
				Name:     "IndexObject",
				Type:     "TPanel",
				HasIndex: true,
				Index:    123,
			}},
			{Name: "Left", Value: dfm.Int(123)},
			{Name: "Top", Value: dfm.Int(-123)},
			{Name: "Scale", Value: dfm.Float(1.0)},
			{Name: "F.G", Value: dfm.Float(-123.1875)},
			{Name: "Precise", Value: dfm.Float(39043.36641510417)},
			{Name: "Huge", Value: dfm.Float(1.000000040918479e35)},
			{Name: "Decimal", Value: dfm.Float(1e+15)},
			{Name: "UseE", Value: dfm.Float(1e+16)},
			{Name: "Not.A.Number", Value: dfm.Float(math.NaN())},
			{Name: "Infinity", Value: dfm.Float(math.Inf(+1))},
			{Name: "NegativeInfinity", Value: dfm.Float(math.Inf(-1))},
			{Name: "EmptyString", Value: dfm.String("")},
			{Name: "S", Value: dfm.String("string")},
			{Name: "Unicode1Line", Value: dfm.String(strings.Repeat("ä", 64))},
			{Name: "Unicode2Lines", Value: dfm.String(strings.Repeat("ä", 65))},
			{Name: "OneLine", Value: dfm.String(strings.Repeat("x", 64))},
			{Name: "TwoLines", Value: dfm.String(strings.Repeat("x", 65))},
			{Name: "Quoted", Value: dfm.String("The 'Laser'")},
			{Name: "Control", Value: dfm.String("\t\r\n")},
			{Name: "LongString", Value: dfm.String(`
				A long string with four tabs at the start of each line.
				The string starts with a line break and ends with a line break,
				followed by some more tabs.
			`)},
			{Name: "Yes", Value: dfm.Bool(true)},
			{Name: "No", Value: dfm.Bool(false)},
			{Name: "ID", Value: dfm.Identifier("clColor")},
			{Name: "ID.With.Dots", Value: dfm.Identifier("A.B.C")},
			{Name: "EmptySet", Value: dfm.Set{}},
			{Name: "Left", Value: dfm.Set{
				dfm.Identifier("akLeft"),
			}},
			{Name: "TopLeft", Value: dfm.Set{
				dfm.Identifier("akLeft"),
				dfm.Identifier("akTop"),
			}},
			{Name: "EmptyTuple", Value: dfm.Tuple{}},
			{Name: "One", Value: dfm.Tuple{dfm.Int(1)}},
			{Name: "OneTwo", Value: dfm.Tuple{dfm.Int(1), dfm.Int(2)}},
			{Name: "StringTuple", Value: dfm.Tuple{
				dfm.String(strings.Repeat("a", 10)),
				dfm.String(strings.Repeat("b", 100)),
				dfm.String(strings.Repeat("c", 5)),
			}},
			{Name: "EmptyBytes", Value: dfm.Bytes{}},
			{Name: "OneByte", Value: dfm.Bytes{0xAF}},
			{Name: "ManyBytes", Value: dfm.Bytes(bytes.Repeat([]byte{
				0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}, 10))},
			{Name: "EmptyItems", Value: dfm.Items{}},
			{Name: "OneEmptyItem", Value: dfm.Items{
				[]dfm.Property{},
			}},
			{Name: "TwoEmptyItems", Value: dfm.Items{
				[]dfm.Property{},
				[]dfm.Property{},
			}},
			{Name: "OneItem", Value: dfm.Items{
				[]dfm.Property{
					{Name: "Left", Value: dfm.Int(5)},
				},
			}},
			{Name: "NestedItems", Value: dfm.Items{
				[]dfm.Property{
					{Name: "Nested", Value: dfm.Items{
						[]dfm.Property{
							{Name: "Left", Value: dfm.Int(5)},
						},
					}},
				},
			}},
		},
	}

	want := `object Dialog: TDialog
  inline TSubObject
  end
  inherited Child: TChild
  end
  object IndexObject: TPanel [123]
  end
  Left = 123
  Top = -123
  Scale = 1.000000000000000000
  F.G = -123.187500000000000000
  Precise = 39043.366415104170000000
  Huge = 1.000000040918479E35
  Decimal = 1000000000000000.000000000000000000
  UseE = 1E16
  Not.A.Number = 0.000000000000000000
  Infinity = 0.000000000000000000
  NegativeInfinity = 0.000000000000000000
  EmptyString = ''
  S = 'string'
  Unicode1Line = #228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228
  Unicode2Lines = 
    #228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228#228 +
    #228
  OneLine = 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
  TwoLines = 
    'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' +
    'x'
  Quoted = 'The '#39'Laser'#39
  Control = #9#13#10
  LongString = 
    #10#9#9#9#9'A long string with four tabs at the start of each line.'#10#9#9#9 +
    #9'The string starts with a line break and ends with a line break,' +
    #10#9#9#9#9'followed by some more tabs.'#10#9#9#9
  Yes = True
  No = False
  ID = clColor
  ID.With.Dots = A.B.C
  EmptySet = []
  Left = [akLeft]
  TopLeft = [akLeft, akTop]
  EmptyTuple = ()
  One = (
    1)
  OneTwo = (
    1
    2)
  StringTuple = (
    'aaaaaaaaaa'
    
      'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' +
      'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb'
    'ccccc')
  EmptyBytes = {
    }
  OneByte = {
    AF}
  ManyBytes = {
    0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF
    0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF
    0123456789ABCDEF0123456789ABCDEF}
  EmptyItems = <>
  OneEmptyItem = <
    item
    end>
  TwoEmptyItems = <
    item
    end
    item
    end>
  OneItem = <
    item
      Left = 5
    end>
  NestedItems = <
    item
      Nested = <
        item
          Left = 5
        end>
    end>
end
`

	// Go literal break lines with \n where Delphi uses \r\n.
	want = strings.Replace(want, "\n", "\r\n", -1)
	have := obj.String()
	if have != want {
		t.Errorf(
			"wrong DFM printed, want:\n---\n%s\n---\nbut have:\n---\n%s\n---",
			want, have,
		)
	}
}
