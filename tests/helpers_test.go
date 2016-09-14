package confmgr

import (
	"github.com/moensch/confmgr"
	"testing"
)

func TestTypeToString(t *testing.T) {
	testdata := make(map[int]string)

	testdata[0] = "none"
	testdata[1] = "string"
	testdata[3] = "list"
	testdata[5] = "hash"
	testdata[2] = "INVALID"

	for int_type, expected := range testdata {
		t.Logf("Testing data type %d", int_type)
		actual := confmgr.TypeToString(int_type)
		t.Logf("  Expected: %s", expected)
		t.Logf("  Actual  : %s", actual)
		if actual != expected {
			t.Fail()
		}
	}
}
