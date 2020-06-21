/*
Package dfm implements a Delphi DFM (.dfm) file parser, generator and printer.

Use any of these functions to parse a DFM file:

	ParseString(code string)
	ParseBytes(code []byte)
	ParseFile(path string)
	ParseReader(r io.Reader)

They all return a dfm.Object and error.

A DFM file contains one root Object which contains other objects and properties,
forming a tree structure. Properties can be of types (see file dfm.go):

	Int
	Float
	Bool
	String
	Identifier
	Bytes
	Set
	Tuple
	Items

You can maipulate the in-memory tree by replacing its nodes.

To write an Object to a file, use any of these functions:

	Object.Print() []byte
	Object.String() string
	Object.WriteTo(w io.Writer) error

These will create an ASCII or UTF-8 encoded (depending on whether the DFM
contains unicode characters in its identifiers) code file, readable by Delphi.
*/
package dfm
