package dfm_test

import (
	"testing"

	"github.com/gonutz/check"
	"github.com/gonutz/dfm"
)

func TestParseEmptyObject(t *testing.T) {
	parseObject(t,
		`object Dialog: TDialog
end`,
		&dfm.Object{
			Name: "Dialog",
			Type: "TDialog",
		},
	)
}

func TestParseFileWithUTF8Header(t *testing.T) {
	code := append(utf8bom, []byte(`object T end`)...)
	obj, err := dfm.ParseBytes(code)
	check.Eq(t, err, nil)
	check.Eq(t, obj, &dfm.Object{Type: "T"})
}

func TestParseIntProperty(t *testing.T) {
	parseProperties(t, `
  Left = 0
  Top = -11`,
		dfm.Property{Name: "Left", Value: dfm.Int(0)},
		dfm.Property{Name: "Top", Value: dfm.Int(-11)},
	)
}

func TestParsePropertyNameWithDot(t *testing.T) {
	parseProperties(t,
		"VertScrollBar.Position = 3231",
		dfm.Property{Name: "VertScrollBar.Position", Value: dfm.Int(3231)},
	)
}

func TestParseBooleanProperty(t *testing.T) {
	parseProperties(t, `
  Old = True
  New = False`,
		dfm.Property{Name: "Old", Value: dfm.Bool(true)},
		dfm.Property{Name: "New", Value: dfm.Bool(false)},
	)
}

func TestParseStringProperty(t *testing.T) {
	parseProperties(t,
		"S = 'some string'",
		dfm.Property{Name: "S", Value: dfm.String("some string")},
	)
}

func TestParseStringPropertyWithQuotes(t *testing.T) {
	parseProperties(t,
		"Q = 'quoted '' string'",
		dfm.Property{Name: "Q", Value: dfm.String("quoted ' string")},
	)
}

func TestParseStringPropertyWithCharacters(t *testing.T) {
	parseProperties(t,
		"Unicode = 'Hyphen '#8211' and umlaut '#252' in string'",
		dfm.Property{
			Name:  "Unicode",
			Value: dfm.String("Hyphen – and umlaut ü in string"),
		},
	)
}

func TestParseConcatenatedStringProperty(t *testing.T) {
	parseProperties(t, `
  S = 'This is the first line ' +
      'while this is the second.'`,
		dfm.Property{
			Name:  "S",
			Value: dfm.String("This is the first line while this is the second."),
		},
	)
}

func TestParseIdentifierProperty(t *testing.T) {
	parseProperties(t,
		"Color = clRed",
		dfm.Property{Name: "Color", Value: dfm.Identifier("clRed")},
	)
}

func TestParseIdentifierWithDotsProperty(t *testing.T) {
	parseProperties(t,
		"Event = Device.Action",
		dfm.Property{Name: "Event", Value: dfm.Identifier("Device.Action")},
	)
}

func TestParseSetProperty(t *testing.T) {
	parseProperties(t,
		"Anchors = [akLeft, akTop, akRight]",
		dfm.Property{Name: "Anchors", Value: dfm.Set{
			dfm.Identifier("akLeft"),
			dfm.Identifier("akTop"),
			dfm.Identifier("akRight"),
		}},
	)
}

func TestParseTupleProperty(t *testing.T) {
	parseProperties(t, `
  DesignSize = (
    800
    4800)`,
		dfm.Property{Name: "DesignSize", Value: dfm.Tuple{
			dfm.Int(800),
			dfm.Int(4800),
		}},
	)
}

func TestParseTupleOfStrings(t *testing.T) {
	parseProperties(t, `
  Strings = (
    'a'
    'b'
    
    'bro' +

    'ken')`,
		dfm.Property{Name: "Strings", Value: dfm.Tuple{
			dfm.String("a"),
			dfm.String("b"),
			dfm.String("broken"),
		}},
	)
}

func TestParseNestedObject(t *testing.T) {
	parseObject(t,
		`object A: TA
  object B: TB
    Nested = True
  end
end`,
		&dfm.Object{
			Name: "A",
			Type: "TA",
			Properties: []dfm.Property{
				{
					Name: "B",
					Value: &dfm.Object{
						Name: "B",
						Type: "TB",
						Properties: []dfm.Property{
							{Name: "Nested", Value: dfm.Bool(true)},
						},
					},
				},
			},
		},
	)
}

func TestParseNestedObjects(t *testing.T) {
	parseObject(t,
		`object A: TA
  object B: TB
  end
  object C: TC
  end
end`,
		&dfm.Object{
			Name: "A",
			Type: "TA",
			Properties: []dfm.Property{
				{
					Name: "B",
					Value: &dfm.Object{
						Name:       "B",
						Type:       "TB",
						Properties: []dfm.Property{},
					},
				},
				{
					Name: "C",
					Value: &dfm.Object{
						Name:       "C",
						Type:       "TC",
						Properties: []dfm.Property{},
					},
				},
			},
		},
	)
}

func TestParseBytesProperty(t *testing.T) {
	parseProperties(t, `
  Bitmap.Data = {
    0123
    4567
    89AB
    cdef}
  A=0`,
		dfm.Property{
			Name:  "Bitmap.Data",
			Value: dfm.Bytes{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF},
		},
		dfm.Property{Name: "A", Value: dfm.Int(0)},
	)
}

func TestParseFloatProperties(t *testing.T) {
	parseProperties(t, `
  Scale = 1.5
  Height = -2.5
  Pos = +10.000
  Tiny = 5.5E-10
  Huge = 1.1e20`,
		dfm.Property{Name: "Scale", Value: dfm.Float(1.5)},
		dfm.Property{Name: "Height", Value: dfm.Float(-2.5)},
		dfm.Property{Name: "Pos", Value: dfm.Float(10)},
		dfm.Property{Name: "Tiny", Value: dfm.Float(5.5e-10)},
		dfm.Property{Name: "Huge", Value: dfm.Float(1.1e+20)},
	)
}

func TestParseInheritedObject(t *testing.T) {
	parseObject(t,
		`inherited Dialog: TDialog
end`,
		&dfm.Object{
			Name: "Dialog",
			Type: "TDialog",
			Kind: dfm.Inherited,
		},
	)
}

func TestParseInlineObject(t *testing.T) {
	parseObject(t,
		`inline Dialog: TDialog
end`,
		&dfm.Object{
			Name: "Dialog",
			Type: "TDialog",
			Kind: dfm.Inline,
		},
	)
}

func TestParseEmptyItemList(t *testing.T) {
	parseProperties(t, `
  EmptyList = <>`,
		dfm.Property{Name: "EmptyList", Value: dfm.Items{}},
	)
}

func TestParseOneItemListWithoutProperties(t *testing.T) {
	parseProperties(t, `
  List = <
    item
    end>`,
		dfm.Property{Name: "List", Value: dfm.Items{[]dfm.Property{}}},
	)
}

func TestParseItemLists(t *testing.T) {
	parseProperties(t, `
  EmptyList = <>
  List = <
    item
    end
    item
      Zero = 0
    end
    item
      One = 1
      Two = 2
    end>`,
		dfm.Property{Name: "EmptyList", Value: dfm.Items{}},
		dfm.Property{Name: "List", Value: dfm.Items{
			[]dfm.Property{},
			[]dfm.Property{
				{Name: "Zero", Value: dfm.Int(0)},
			},
			[]dfm.Property{
				{Name: "One", Value: dfm.Int(1)},
				{Name: "Two", Value: dfm.Int(2)},
			},
		}},
	)
}

func TestParseAnonymousObject(t *testing.T) {
	parseObject(t,
		`object TMenuItem
end`,
		&dfm.Object{
			Name: "",
			Type: "TMenuItem",
		},
	)
}

func TestParseObjectIndex(t *testing.T) {
	parseObject(t,
		`object P0: TPanel [0]
end`,
		&dfm.Object{
			Name:     "P0",
			Type:     "TPanel",
			HasIndex: true,
			Index:    0,
		},
	)

	parseObject(t,
		`object P1: TPanel [1]
end`,
		&dfm.Object{
			Name:     "P1",
			Type:     "TPanel",
			HasIndex: true,
			Index:    1,
		},
	)
}

func TestNonASCIIBytesAreInterpretedAsWindowsANSI(t *testing.T) {
	code := []byte("object O A='\xA9' end") // 0xA9 is the ANSI copyright symbol
	obj, err := dfm.ParseBytes(code)
	check.Eq(t, err, nil)
	check.Eq(t, obj, &dfm.Object{Type: "O", Properties: []dfm.Property{
		{Name: "A", Value: dfm.String("©")},
	}})
}

func parseObject(t *testing.T, code string, want *dfm.Object) {
	t.Helper()
	obj, err := dfm.ParseString(code)
	if err != nil {
		t.Fatal(err)
	}
	check.Eq(t, obj, want)
}

func parseProperties(t *testing.T, partialCode string, want ...dfm.Property) {
	t.Helper()
	code := "object Obj: Typ\r\n" + partialCode + "\r\nend"
	obj, err := dfm.ParseString(code)
	if err != nil {
		t.Fatal(err)
	}
	check.Eq(t, obj, &dfm.Object{
		Name:       "Obj",
		Type:       "Typ",
		Properties: want,
	})
}
