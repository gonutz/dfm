package dfm_test

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/gonutz/dfm"
)

func TestPrintDFM(t *testing.T) {
	// Things to test:
	// - if Object.Name is empty, only print the type.
	// - All Kinds of Object.
	// - All implementations of PropertyValue.
	// - Print Float with 18 digits after the point.
	// - Strings with quotes and non-ASCII characters.
	// - Long Strings should span multiple lines.
	// - Identifiers with dots, both as property names and values.
	// - Floats NAN and +-INF are invalid in DFMs, default to 0.
	obj := dfm.Object{
		Name: "Dialog",
		Type: "TDialog",
		Properties: []dfm.Property{
			dfm.Property{
				Value: dfm.Object{
					Type: "TSubObject",
					Kind: dfm.Inline,
				},
			},
			dfm.Property{
				Value: dfm.Object{
					Name: "Child",
					Type: "TChild",
					Kind: dfm.Inherited,
				},
			},
			dfm.Property{Name: "Left", Value: dfm.Int(123)},
			dfm.Property{Name: "Top", Value: dfm.Int(-123)},
			dfm.Property{Name: "Scale", Value: dfm.Float(1.0)},
			dfm.Property{Name: "F.G", Value: dfm.Float(-123.1875)},
			dfm.Property{Name: "Not.A.Number", Value: dfm.Float(math.NaN())},
			dfm.Property{Name: "Infinity", Value: dfm.Float(math.Inf(+1))},
			dfm.Property{Name: "NegativeInfinity", Value: dfm.Float(math.Inf(-1))},
			dfm.Property{Name: "EmptyString", Value: dfm.String("")},
			dfm.Property{Name: "S", Value: dfm.String("string")},
			dfm.Property{Name: "Quoted", Value: dfm.String("The 'Laser'")},
			dfm.Property{Name: "NonASCII", Value: dfm.String("\t\r\n")},
			dfm.Property{Name: "LongString", Value: dfm.String(`
				A long string with four tabs at the start of each line.
				The string starts with a line break and ends with a line break,
				followed by some more tabs.
			`)},
			dfm.Property{Name: "Yes", Value: dfm.Bool(true)},
			dfm.Property{Name: "No", Value: dfm.Bool(false)},
			dfm.Property{Name: "ID", Value: dfm.Identifier("clColor")},
			dfm.Property{Name: "ID.With.Dots", Value: dfm.Identifier("A.B.C")},
			dfm.Property{Name: "EmptySet", Value: dfm.Set{}},
			dfm.Property{Name: "Left", Value: dfm.Set{
				dfm.Identifier("akLeft"),
			}},
			dfm.Property{Name: "TopLeft", Value: dfm.Set{
				dfm.Identifier("akLeft"),
				dfm.Identifier("akTop"),
			}},
			dfm.Property{Name: "EmptyTuple", Value: dfm.Tuple{}},
			dfm.Property{Name: "One", Value: dfm.Tuple{
				dfm.Int(1),
			}},
			dfm.Property{Name: "OneTwo", Value: dfm.Tuple{
				dfm.Int(1),
				dfm.Int(2),
			}},
			dfm.Property{Name: "EmptyBytes", Value: dfm.Bytes{}},
			dfm.Property{Name: "OneByte", Value: dfm.Bytes{0xAF}},
			dfm.Property{
				Name: "ManyBytes",
				Value: dfm.Bytes(bytes.Repeat([]byte{
					0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF,
				}, 10)),
			},
			dfm.Property{Name: "EmptyItems", Value: dfm.Items{}},
			dfm.Property{Name: "OneEmptyItem", Value: dfm.Items{
				[]dfm.Property{},
			}},
			dfm.Property{Name: "TwoEmptyItems", Value: dfm.Items{
				[]dfm.Property{},
				[]dfm.Property{},
			}},
			dfm.Property{Name: "OneItem", Value: dfm.Items{
				[]dfm.Property{
					dfm.Property{Name: "Left", Value: dfm.Int(5)},
				},
			}},
			dfm.Property{Name: "NestedItems", Value: dfm.Items{
				[]dfm.Property{
					dfm.Property{Name: "Nested", Value: dfm.Items{
						[]dfm.Property{
							dfm.Property{Name: "Left", Value: dfm.Int(5)},
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
  Left = 123
  Top = -123
  Scale = 1.000000000000000000
  F.G = -123.187500000000000000
  Not.A.Number = 0.000000000000000000
  Infinity = 0.000000000000000000
  NegativeInfinity = 0.000000000000000000
  EmptyString = ''
  S = 'string'
  Quoted = 'The '#39'Laser'#39
  NonASCII = #9#13#10
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
