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
