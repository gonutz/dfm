Package `dfm` is a library, written in Go, to parse, pretty-print and/or generate Delphi's DFM files.

[See the Godoc documentation for details of the API.](https://godoc.org/github.com/gonutz/dfm)

You can parse a DFM file with either of these functions:

	dfm.ParseString(code string)
	dfm.ParseBytes(code []byte)
	dfm.ParseFile(path string)
	dfm.ParseReader(r io.Reader)

which all return a pointer to a `dfm.Object`.

You can manipulate Objects in memory, either Objects that were parsed from an existing DFM file or you can create a new Object from scratch. These can be written back to file to be used in Delphi.

To generate code from an Object you can call one of these functions:

	dfm.Object.Print() []byte
	dfm.Object.String() string
	dfm.Object.WriteTo(w io.Writer) error

The generated code is formatted exactly like RAD Studio XE4 formats it. It will almost always match the file byte for byte. The only known difference is that floating point numbers might appear slightly different. Their value will be the same but some files use 15 decimal points while most use 18. This library was tweaked to match all of the available test DFM files exactly, only two of them differ in this way. If you encounter any significant differenes, please provide the sample DFM in a [Github issue](https://github.com/gonutz/dfm/issues).
The output DFMs will be encoded in ASCII, except if any of the identifiers use non-ASCII characters, in that case the code is encoded as UTF-8 and starts with the UTF-8 byte order mark. This matches RAD Studio behavior.

This library was tested against 600 DFM files from both the RAD Studio sources and production code from the company I work at. All files are parsed correctly and printed back to produce the exact same file as was input, except from two minor issues (see above). If you encounter any problems, please write a [Github issue](https://github.com/gonutz/dfm/issues).
