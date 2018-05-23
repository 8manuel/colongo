package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/8manuel/colongo/colon"
	"github.com/stellar/go/keypair"
)

// getAssetKeypairs generates a pair of keypairs (seed+address), issuing and distribution for issuing an asset
func getAssetKeypairs() (pairIss, pairDis *keypair.Full, err error) {
	// generate a [32]byte seed for the issuing account
	bytesIss := []byte("issuing account byteseed that can be printed")
	byteSeedIss := [32]byte{}
	for i := 0; i < 32; i++ {
		byteSeedIss[i] = bytesIss[i]
	}
	// generate the keypair from the seed
	pairIss, err = keypair.FromRawSeed(byteSeedIss)
	if err != nil {
		return pairIss, pairDis, err
	}

	// generate a [32]byte seed for the distribution account
	bytesDis := []byte("distribution account byteseed that can be printed")
	byteSeedDis := [32]byte{}
	for i := 0; i < 32; i++ {
		byteSeedDis[i] = bytesDis[i]
	}
	// generate the keypair from the seed
	pairDis, err = keypair.FromRawSeed(byteSeedDis)
	if err != nil {
		return pairIss, pairDis, err
	}

	return pairIss, pairDis, err
}

func fundAddress(addr string) (err error) {
	// ask balance to the faucet bot
	log.Printf("Requesting funding for Address %s\n", addr)
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + addr)
	if err != nil {
		return err
	}
	// check the response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		fmt.Println(body)
	}
	return err
}

func balAddress(addr string) (err error) {
	account, err := colon.MLoadAccount(addr)
	if err == nil {
		for _, balance := range account.Balances {
			fmt.Println("Addr", addr, balance)
		}
	}
	return err
}

func TestAssetAddrs(t *testing.T) {
	pairIss, pairDis, err := getAssetKeypairs()
	if err != nil {
		t.Error(err)
	}
	// get and print the seed and the address
	seedIss, addrIss := pairIss.Seed(), pairIss.Address()
	log.Printf("Issuing keypair Seed %s, Address %s\n", seedIss, addrIss)

	// get and print the seed and the address
	seedDis, addrDis := pairDis.Seed(), pairDis.Address()
	log.Printf("Distribution keypair Seed %s, Address %s\n", seedDis, addrDis)
}

func TestAssetFund(t *testing.T) {
	// get the issuing and distribution keypairs
	pairIss, pairDis, err := getAssetKeypairs()
	if err != nil {
		t.Error(err)
	}
	// fund both addresses
	if err = fundAddress(pairIss.Address()); err != nil {
		t.Error(err)
	}
	if err = fundAddress(pairDis.Address()); err != nil {
		t.Error(err)
	}
}

func TestAssetBal(t *testing.T) {
	// get the issuing and distribution keypairs
	pairIss, pairDis, err := getAssetKeypairs()
	if err != nil {
		t.Error(err)
	}
	// check both addresses balance
	if err = balAddress(pairIss.Address()); err != nil {
		t.Error(err)
	}
	if err = balAddress(pairDis.Address()); err != nil {
		t.Error(err)
	}

}

func TestAssetTrust(t *testing.T) {
	// get the issuing and distribution keypairs
	pairIss, pairDis, err := getAssetKeypairs()
	if err != nil {
		t.Error(err)
	}
	if err = colon.MTransTrust(pairDis, "VEF", pairIss.Address(), "1500", true); err != nil {
		t.Error(err)
	}
}

func TestTransPayXLM(t *testing.T) {
	// get the issuing and distribution keypairs
	pairIss, pairDis, err := getAssetKeypairs()
	if err != nil {
		t.Error(err)
	}
	_ = pairIss
	if err = colon.MTransPayment(pairDis, "GBUAFDIXDT4EOPLAJN7CWVXSHRN3KH2ASKRNMCIIIMXOO4QYWFIHMBEG", "", "0.1", true); err != nil {
		t.Error(err)
	}
}

func TestTransPayAsset(t *testing.T) {
	// get the issuing and distribution keypairs
	pairIss, pairDis, err := getAssetKeypairs()
	if err != nil {
		t.Error(err)
	}
	if err = colon.MTransPayment(pairIss, pairDis.Address(), "VEF", "1500", true); err != nil {
		t.Error(err)
	}
}

func TestTransAllowTrust(t *testing.T) {
	// get the issuing and distribution keypairs
	pairIss, pairDis, err := getAssetKeypairs()
	if err != nil {
		t.Error(err)
	}
	if err = colon.MAllowTrust(pairIss, "VEF", pairDis.Address(), true, false); err != nil {
		t.Error(err)
	}
}

func TestTransSetOptions(t *testing.T) {
	// get the issuing and distribution keypairs
	pairIss, _, err := getAssetKeypairs()
	if err != nil {
		t.Error(err)
	}
	//if err = colon.MSetOptions(pairIss, map[string]interface{}{"SetFlags": uint32(0x01 | 0x02 | 0x04)}); err != nil { // authReq+authRev+authImm
	//if err = colon.MSetOptions(pairIss, map[string]interface{}{"ClearFlags": uint32(0x02 | 0x04)}); err != nil { // authRev+authImm
	if err = colon.MSetOptions(pairIss, map[string]interface{}{"HomeDomain": "subdomain.domain.com"}); err != nil { // authRev+authImm
		t.Error(err)
	}
}
