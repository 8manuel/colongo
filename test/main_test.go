// To execute all test functions, type in the terminal
//  go test
// To execute an specific test function, like TestAddrBalance, type in the terminal
//  go test -run TestAddrBalance
//
// This code is public in github and I made it to help other people to learn Stellar+golang.
// In package colon there are some functions that interact with the Stellar libraries and help you to do the drills.
//
// NOTE there is a very important parameter "flbBaseSeed" that it is used to deterministic generate keypairs.
// (deterministic means that the same input produces always the same output result)
// In colon.DeterministicKeypair those the 24bytes plus the accName are used to construct deterministic a keypair.
// If you get this code from github if this parameter is not changed the keypairs that you can generate may already exist.
// Therefore I advise to change is when you run the tests. The way to change is to do in your code or to do in the test command line like this
//  go test -run TestAddrBalance flgBaseSeed=MYCUSTOMBASESEED
//
package test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/8manuel/colongo/colon"
)

var flgAmt *int
var flgAddr *string

func TestMain(m *testing.M) {

	var err error
	// load flags
	flgBaseSeed := flag.String("flgBaseSeed", "BaseDrillSeedStr20180522", "used to generate deterministic keypairs for the tests")
	flgAmt = flag.Int("flgAmt", 1000, "Amount")
	flgAddr = flag.String("flgAddr", "GBIYBTHFAOEZNBVDFHAAQWD25EG2CVXCC4PQ333PIUQGRVZN5MJEZRHO", "Address")
	flag.Parse()
	_, _, _ = err, flgAmt, flgAddr
	fmt.Println("running with flags", "flgBaseSeed", *flgBaseSeed, "flgAmt", *flgAmt, "flgAddr", *flgAddr, "\n")

	// set the base seed
	colon.Init(*flgBaseSeed)

	// execute the rest of tests, m.Run() executes tests in the following order:
	// for each file test/*.go file sorted alphabetically
	//	 for each Test* function sorted by line number appearance
	//     execute function
	v := m.Run()

	// exit
	os.Exit(v)
}
