package gcpug

import (
	"github.com/sinmetal/gaego_unittest_util/aetestutil"

	"testing"
)

func TestSpinUpFirst(t *testing.T) {
	_, _, err := aetestutil.SpinUp()
	if err != nil {
		t.Fatal(err)
	}
}
