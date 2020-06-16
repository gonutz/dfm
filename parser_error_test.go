package dfm_test

import (
	"testing"

	"github.com/gonutz/dfm"
)

func TestBinaryDFMfilesAreNotImplemented(t *testing.T) {
	// Binary DFM files start with byte 255. Parsing should stop right there.
	parseError(t, string([]byte{0xFF}), "dfm.Parse: binary DFM files are not supported")
}

func parseError(t *testing.T, code string, wantMsg string) {
	t.Helper()
	_, err := dfm.Parse(code)
	if err == nil {
		t.Fatal("no error")
	}
	haveMsg := err.Error()
	if haveMsg != wantMsg {
		t.Errorf("wrong message, want:\n%s\nbut have:\n%s", wantMsg, haveMsg)
	}
}
