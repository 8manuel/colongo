package test

import (
	"testing"
)

// The objective of this drill is to practice multisignature; it is required to have completed drill0 because the same accounts and assets are used.
// The asset is issued by account A that acts as issuer, the code is VEF.
// Just remember that for the drills we are using testnet and the VEF value is the same as monopoly notes (nothing, well in this case the real VEF value is near to 0).

// TestDrill1Exchange; transacion with two operations, first accB sends 5XLM to accA, second accA sends 15VEF to accB
func TestDrill1Exchange(t *testing.T) {
	// accB generates transaction tx with: op1 5XLM from accB to accA + op2 15VEF from accA to accB
	// sign transaction only with accB seed
	// send transaction to horizon: as accA has not signed "tx_failed", operations:["op_success", "op_bad_auth"]

	// now sign the same tx with accB and accA seed
	// send transaction to horizon; now it should work

	// now sign the same tx with accB and accA and accC seed
	// send transaction to horizon; now it should work even if in the documents say that it fails when there are more transactions than required (I suppose is because accC has no operation and is ignored)
}

// TestDrill1Multi; for accC authorize signatures of accA with weight 1 and accB with weight 1; also set accC threshold for mid to 2 and high to 3.
func TestDrill1Multi(t *testing.T) {
	// accC authorize signatures of accA with weight 1 and accB with weight 1; also set accC threshold for mid to 2 and high to 3

	// send 1XLM from accC to accB; sign only with accC, gives transaction:"tx_failed", operations:["op_bad_auth"]

	// send 1XLM from accC to accB; sign with accA, accB and accC, gives transaction:"tx_bad_auth_extra"

	// send 1XLM from accC to accB; sign with accB and accC, now it works

}
