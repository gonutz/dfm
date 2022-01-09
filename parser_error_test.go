package dfm_test

import (
	"testing"

	"github.com/gonutz/check"
	"github.com/gonutz/dfm"
)

func TestBinaryDFMfilesAreNotImplemented(t *testing.T) {
	// Binary DFM files start with byte 0xFF. Parsing should stop right there.
	_, err := dfm.ParseBytes([]byte{0xFF})
	check.Eq(t, err.Error(), "dfm.Parse: binary DFM files are not supported")
}

func TestInvalidIntegerInItems(t *testing.T) {
	_, err := dfm.ParseString(`object O: TO
  List = <
    item
	  Value = 123d
    end>
end`)
	check.Neq(t, err, nil)
}

func TestMissingCommaInSet(t *testing.T) {
	_, err := dfm.ParseString(`object O: TO
	Set = [123 456]
end`)
	check.Neq(t, err, nil)
}

func TestInvalidIntegerInSet(t *testing.T) {
	_, err := dfm.ParseString(`object O: TO
	Set = [123d]
end`)
	check.Neq(t, err, nil)
}
