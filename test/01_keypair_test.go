package test

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/8manuel/colongo/colon"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

//
// addresses test using keypair package
//

func TestAddrGenRandom1(t *testing.T) {

	// generate a random keypair
	pair, err := keypair.Random()
	if err != nil {
		t.Error(err)
	}
	// print the seed and address
	log.Println(pair.Seed(), pair.Address())
}

func TestAddrGenRandom2(t *testing.T) {

	// generate a random keypair
	var pair *keypair.Full
	var err error
	if pair, err = keypair.Random(); err != nil {
		t.Error(err)
	}
	// get the seed and the address
	seed, addr := pair.Seed(), pair.Address()
	// print the result
	log.Printf("Seed %s, len %d , Address %s\n", seed, len(seed), addr)

	// convert the seed to []byte
	bs, err := colon.MSeed2Bytes(seed)
	log.Printf("Extracted from seed %s, this %X\n", seed, bs)
}

func TestAddrGenDet0(t *testing.T) {
	// generate a keypair using colon package function
	pair := colon.DeterministicKeypair("Det0")
	if pair == nil {
		t.Error(errors.New("failed DeterministicKeypair"))
		return
	}
	// get and print the seed and the address
	seed, addr := pair.Seed(), pair.Address()
	log.Printf("Keypair Seed %s, Address %s\n", seed, addr)

	// convert the seed to []byte
	btyeSeed, err := colon.MSeed2Bytes(seed)
	if err != nil {
		t.Error(err)
	}
	log.Printf("Extracted from seed %s, this \"%s\"\n", seed, btyeSeed)
}

func TestAddrGenDet1(t *testing.T) {
	// generate a [32]byte seed
	bytes := []byte("no se como hacer esto asi que me lo invento")
	//bytes = []byte{110, 111, 32, 67, 68, 69, 70, 71, 64, 65, 66, 67, 68, 69, 70, 71, 64, 65, 66, 67, 68, 69, 70, 71, 64, 65, 66, 67, 68, 69, 70, 71}
	byteSeed := [32]byte{}
	for i := 0; i < 32; i++ {
		byteSeed[i] = bytes[i]
	}
	log.Printf("Using byteSeed \"%s\", generate a keypair\n", byteSeed)

	// generate the keypair from the seed
	pair, err := keypair.FromRawSeed(byteSeed)
	if err != nil {
		t.Error(err)
	}

	// get and print the seed and the address
	seed, addr := pair.Seed(), pair.Address()
	log.Printf("Keypair Seed %s, Address %s\n", seed, addr)

	// convert the seed to []byte
	btyeSeed, err := colon.MSeed2Bytes(seed)
	log.Printf("Extracted from seed %s, this \"%s\"\n", seed, btyeSeed)
}

func TestAddrParse(t *testing.T) {
	// parse a seed SDNYODGEMGKGIBNCR6C6XYQ7LUH5CIL2MNNIDTQQPWO6XNTIAVRHF43P into a keypair with addr GAMDGSB4NNTY4GIQMSQVSDMDNKAVDQTHRS7HSL5JJJRRBFKWHPJZHA7K
	pair, err := keypair.Parse("SDNYODGEMGKGIBNCR6C6XYQ7LUH5CIL2MNNIDTQQPWO6XNTIAVRHF43P")
	if err != nil {
		t.Error(err)
	}
	//	log.Printf("Seed %s, Address %s\n", pair.Seed(), pair.Address())
	log.Printf("Seed %s, Address %s\n", pair, pair.Address())
}

func TestAddrBalance(t *testing.T) {
	account, err := horizon.DefaultTestNetClient.LoadAccount(*flgAddr)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Balances for account:", *flgAddr)
	for _, balance := range account.Balances {
		log.Println(balance)
	}
}
