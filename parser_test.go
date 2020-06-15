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
		dfm.Object{
			Name: "Dialog",
			Type: "TDialog",
		},
	)
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

func TestParseNestedObject(t *testing.T) {
	parseObject(t,
		`object A: TA
  object B: TB
    Nested = True
  end
end`,
		dfm.Object{
			Name: "A",
			Type: "TA",
			Properties: []dfm.Property{
				dfm.Property{
					Name: "B",
					Value: dfm.Object{
						Name: "B",
						Type: "TB",
						Properties: []dfm.Property{
							dfm.Property{Name: "Nested", Value: dfm.Bool(true)},
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
		dfm.Object{
			Name: "A",
			Type: "TA",
			Properties: []dfm.Property{
				dfm.Property{
					Name: "B",
					Value: dfm.Object{
						Name:       "B",
						Type:       "TB",
						Properties: []dfm.Property{},
					},
				},
				dfm.Property{
					Name: "C",
					Value: dfm.Object{
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
  Pos = +10.000`,
		dfm.Property{Name: "Scale", Value: dfm.Float(1.5)},
		dfm.Property{Name: "Height", Value: dfm.Float(-2.5)},
		dfm.Property{Name: "Pos", Value: dfm.Float(10)},
	)
}

func parseObject(t *testing.T, code string, want dfm.Object) {
	t.Helper()
	obj, err := dfm.Parse(code)
	if err != nil {
		t.Fatal(err)
	}
	check.Eq(t, obj, want)
}

func parseProperties(t *testing.T, partialCode string, want ...dfm.Property) {
	t.Helper()
	code := "object Obj: Typ\r\n" + partialCode + "\r\nend"
	obj, err := dfm.Parse(code)
	if err != nil {
		t.Fatal(err)
	}
	check.Eq(t, obj, dfm.Object{
		Name:       "Obj",
		Type:       "Typ",
		Properties: want,
	})
}
