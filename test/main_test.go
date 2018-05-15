package test

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	var err error
	// load flags, with "flags" all flag are overriden (except logLevel)
	flgAmt := flag.Int("flgAddr", 1000, "Amount")
	flgAddr := flag.String("flgAddr", "SDNYODGEMGKGIBNCR6C6XYQ7LUH5CIL2MNNIDTQQPWO6XNTIAVRHF43P ", "Address")
	flag.Parse()

	// execute the rest of tests, m.Run() executes tests in the following order:
	// for each file test/*.go file sorted alphabetically
	//	 for each Test* function sorted by line number appearance
	//     execute function
	v := m.Run()

	// exit
	os.Exit(v)
}
