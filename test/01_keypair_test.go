package test

import (
	"encoding/base32"
	"log"
	"testing"

	"github.com/stellar/go/keypair"
)

//
// addresses test using keypair package
//

func TestAddrGenRandom1(t *testing.T) {

	// generate a random keypair
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}
	// print the seed and address
	log.Println(pair.Seed(), pair.Address())
}

func TestAddrGenRandom2(t *testing.T) {

	// generate a random keypair
	var pair *keypair.Full
	var err error
	if pair, err = keypair.Random(); err != nil {
		log.Fatal(err)
	}
	// get the seed and the address
	seed, addr := pair.Seed(), pair.Address()
	// print the result
	log.Printf("Seed %s, len %d , Address %s\n", seed, len(seed), addr)
}

func TestAddrGenDet1(t *testing.T) {
	// generate a [32]byte seed
	bytes := []byte("no se como hacer esto asi que me lo invento")
	byteSeed := [32]byte{}
	for i := 0; i < 32; i++ {
		byteSeed[i] = bytes[i]
	}
	// generate the keypair from the seed
	pair, err := keypair.FromRawSeed(byteSeed)
	if err != nil {
		log.Fatal(err)
	}

	// get and print the seed and the address
	seed, addr := pair.Seed(), pair.Address()
	log.Printf("Seed %s, Address %s\n", seed, addr)

	// convert the seed to []byte
	data, err := base32.StdEncoding.DecodeString(pair.Seed())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Seed \"%s\", byteseed \"%s\"\n", byteSeed, data[1:33])
}

func TestAddrParse(t *testing.T) {
	// parse a seed SDNYODGEMGKGIBNCR6C6XYQ7LUH5CIL2MNNIDTQQPWO6XNTIAVRHF43P into a keypair with addr GAMDGSB4NNTY4GIQMSQVSDMDNKAVDQTHRS7HSL5JJJRRBFKWHPJZHA7K
	pair, err := keypair.Parse("SDNYODGEMGKGIBNCR6C6XYQ7LUH5CIL2MNNIDTQQPWO6XNTIAVRHF43P")
	if err != nil {
		log.Fatal(err)
	}
	//	log.Printf("Seed %s, Address %s\n", pair.Seed(), pair.Address())
	log.Printf("Seed %s, Address %s\n", pair, pair.Address())
}
