package gcpug

import (
	"github.com/sinmetal/gaego_unittest_util/aetestutil"

	"testing"
)

func TestSpinDownLast(t *testing.T) {
	err := aetestutil.SpinDown()
	if err != nil {
		t.Fatal(err)
	}
}
