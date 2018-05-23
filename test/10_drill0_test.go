package test

import (
	"testing"
)

// The objective of this drill is to practice issuing assets from issuer to distributor and to play with trust lines (distributor and holder) and allow trust (issuer).
// The asset is issued by account A that acts as issuer, the code is VEF.
// Just remember that for the drills we are using testnet and the VEF value is the same as monopoly notes (nothing, well in this case the real VEF value is near to 0).

// TestDrill0FundAB creates and funds accounts A and B
func TestDrill0FundAB(t *testing.T) {
	// generate accounts keypairs for names A, B

	// create&fund the accounts; for acc A use friendbot
}

// TestDrill0IssuerSet sets auth required and revocable flags for account A
func TestDrill0IssuerSet(t *testing.T) {
	// set auth required and revocable flags for account A
}

// TestDrill0IssDis; account A creates VEF asset; then send 100 VEF asset from account A (issuer) to account B (distributor); then acc B sends 4 VEF to acc A
func TestDrill0IssDis(t *testing.T) {
	// send 100 VEF asset from account A (issuer) to account B (distributor); as there is no trust gives transaction:"tx_failed", operations:["op_no_trust"]

	// create a trustline from account B to A for 500 VEF

	// send 100 VEF asset from account A (issuer) to account B (distributor); as B is not authorized by A gives transaction:"tx_failed", operations:["op_not_authorized"]

	// account A allows trust to B

	// send 100 VEF asset from account A (issuer) to account B (distributor); now it should work

	// account A revokes trust to B

	// send 4 VEF asset from B to A; as B is not authorized by A gives transaction: "tx_failed", operations:["op_src_not_authorized"]

	// account A allows trust to B

	// send 4 VEF asset from B to A; now it should work, as acc A is the issuer the balance of acc A in VEF will be zero, and acc B balance will be 96 VEF

	// change trustline from account B to A for 50 VEF; as acc B holds 98 it gives transaction: "tx_failed", operations:["op_not_authorized"]

}

// TestDrill0Fund creates and funds account C
func TestDrill0FundC(t *testing.T) {
	// generate accounts keypair for name C

	// send 100XLM from account B to C; as acc C does not exist gives transaction:"tx_failed", operations:["op_no_destination"]

	// create&fund the account; for acc C use friendbot
}

// TestDrill0DisHol; account B (distributor) sends 10 VEF (issuer is acc A) to account C (holder), then acc C sends 5 VEF to acc A
func TestDrill0DisHol(t *testing.T) {
	// send 10 VEF asset from account B (distributor) to account C (holder); as there is no trust gives transaction:"tx_failed", operations:["op_no_trust"]

	// create a trustline from account C to A for 500 VEF

	// send 10 VEF asset from account B (distributor) to account C (holder); as C is not authorized by A gives transaction:"tx_failed", operations:["op_not_authorized"]

	// account A allows trust to C

	// send 10 VEF asset from account B (distributor) to account C (holder); now it should work acc B balance 86 VEF, acc C balance 10 VEF

	// account A revokes trust to C

	// send 5 VEF asset from C to A; as C is not authorized by A gives "transaction": "tx_failed", operations:["op_src_not_authorized"]

	// account A allows trust to C

	// send 5 VEF asset from C to A; now it should work, as acc A is the issuer the balance of acc A in VEF will be zero, and acc C balance will be 5 VEF
}

// TestDrill0IssHol; send 100 VEF asset from account A (issuer) to account C (holder).
// This demonstrates that issuer can send directly to whatever account that has created a trustline with him (in other words, there can be several distributors)
func TestDrill0IssHol(t *testing.T) {
	// send 100 VEF asset from account A (issuer) to account C (distributor); it should work
}

// TestDrill0IssEURb; account A is the issuer of the EUR asset; then send 50 EUR asset from account A (issuer) to account B (distributor).
// Then account A removes trust from B, now both assets (VEF and EUR)9 should be blocked.
func TestDrill0IssEURb(t *testing.T) {
	// create a trustline from account B to A for 200 EUR

	// send 50 EUR asset from account A (issuer) to account B (distributor); it should work

	// account A revokes trust to B

	// send 2 VEF asset from B to A; as B is not authorized by A gives transaction: "tx_failed", operations:["op_src_not_authorized"]

	// send 1 EUR asset from B to A; as B is not authorized by A gives transaction: "tx_failed", operations:["op_src_not_authorized"]
}
