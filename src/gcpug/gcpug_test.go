package gcpug

import (
	"fmt"
	"github.com/sinmetal/gaego_unittest_util/aetestutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("TestMain Start")

	_, _, err := aetestutil.SpinUp()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	status := m.Run()
	if status != 0 {
		fmt.Println("Test Run Error")
	}

	err = aetestutil.SpinDown()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("TestMain End")
}
