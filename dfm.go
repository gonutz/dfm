package dfm

// Parse expects the code to start with an object. The first object in the given
// code is parsed, if there are more, they are ignored. A DFM file typically has
// one top-level object defined in it. It might contain child objects however.
func Parse(code string) (Object, error) {
	return newParser(code).parseObject()
}

// Object can be a TPanel, TLabel, TForm, a sub-class of these or any other
// graphical element that can be defined in Delphi. It contains a list of
// properties, which can include child objects as well.
type Object struct {
	Name       string
	Type       string
	Properties []Property
}

// Property is what is contained in an Object. Possible types are Int, Float,
// Bool, String, Identifier, Set, Tuple, Bytes and Object. Except for Object,
// these will appear in the DFM file as:
//
//     <name> = <value>
//
// where Name can contain dot, e.g. Font.Height.
// In case the Value is an Object, the Name is a copy of the Object.Name.
type Property struct {
	Name  string
	Value PropertyValue
}

// PropertyValue tags types that can be used for a Property.Value.
type PropertyValue interface {
	isPropertyValue()
}

// Int is a decimal integer literal.
type Int int

// Float is a floating point number, represented using a decimal point, i.e.
// rather than using the scientific notations like 2.3E-4 it is written like a
// human, e.g. 0.000234.
type Float float64

// Bool is either True or False.
type Bool bool

// String is a UTF-8 string without enclosing quotes and with quoted quotes
// unquoted. In Delphi we write
//
//     'a ''quoted'' string like this'#13#10
//
// but this string literal in Go is really this:
//
//     "a 'quoted' string like this\r\n"
type String string

// Identifier is a constant like clYellow, poMainFormCenter or FormResize.
type Identifier string

// Set is a set of flags in brackets like
//
//     [akLeft, akTop, akRight]
type Set []PropertyValue

// Tuple is a tuple of values in parentheses like
//
//     (123 456 789)
type Tuple []PropertyValue

// Bytes is a in list of hexadecimal data in braces like
//
//     { FFAC2938AA991234A }
type Bytes []byte

func (Object) isPropertyValue()     {}
func (Int) isPropertyValue()        {}
func (Float) isPropertyValue()      {}
func (Bool) isPropertyValue()       {}
func (String) isPropertyValue()     {}
func (Identifier) isPropertyValue() {}
func (Set) isPropertyValue()        {}
func (Tuple) isPropertyValue()      {}
func (Bytes) isPropertyValue()      {}
