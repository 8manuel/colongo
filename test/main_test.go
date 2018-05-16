// To execute all test functions, type in the terminal
//  go test
// To execute an specific test function, like TestAddrBalance, type in the terminal
//  go test -run TestAddrBalance

package test

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var flgAmt *int
var flgAddr *string

func TestMain(m *testing.M) {

	var err error
	// load flags
	flgAmt = flag.Int("flgAmt", 1000, "Amount")
	flgAddr = flag.String("flgAddr", "GBIYBTHFAOEZNBVDFHAAQWD25EG2CVXCC4PQ333PIUQGRVZN5MJEZRHO", "Address")
	flag.Parse()
	_, _, _ = err, flgAmt, flgAddr
	fmt.Println("running with flags", "flaAmt", *flgAmt, "flgAddr", *flgAddr, "\n")

	// execute the rest of tests, m.Run() executes tests in the following order:
	// for each file test/*.go file sorted alphabetically
	//	 for each Test* function sorted by line number appearance
	//     execute function
	v := m.Run()

	// exit
	os.Exit(v)
}
