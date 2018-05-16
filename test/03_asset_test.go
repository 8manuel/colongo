package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/8manuel/colongo/colon"
	"github.com/stellar/go/clients/horizon"
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

func balXLMAddress(addr string) (err error) {

	account, err := horizon.DefaultTestNetClient.LoadAccount(addr)
	if err != nil {
		return err
	}

	for _, balance := range account.Balances {
		fmt.Println("Balances for account:", addr, balance)
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

func TestAssetBalXLM(t *testing.T) {
	// get the issuing and distribution keypairs
	pairIss, pairDis, err := getAssetKeypairs()
	if err != nil {
		t.Error(err)
	}
	// check both addresses balance
	if err = balXLMAddress(pairIss.Address()); err != nil {
		t.Error(err)
	}
	if err = balXLMAddress(pairDis.Address()); err != nil {
		t.Error(err)
	}

}

func TestAssetTrust(t *testing.T) {
	// get the issuing and distribution keypairs
	pairIss, pairDis, err := getAssetKeypairs()
	if err != nil {
		t.Error(err)
	}
	if err = colon.MTransTrust(pairDis, "VEF", pairIss.Address(), "1000"); err != nil {
		t.Error(err)
	}
}

func TestTransPay(t *testing.T) {
	// get the issuing and distribution keypairs
	pairIss, pairDis, err := getAssetKeypairs()
	if err != nil {
		t.Error(err)
	}
	if err = colon.MTransPayment(pairDis, pairIss.Address(), "10", true); err != nil {
		t.Error(err)
	}

}
