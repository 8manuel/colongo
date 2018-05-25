package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/8manuel/colongo/colon"
	"github.com/stellar/go/build"
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

func TestTransAllowTrust1(t *testing.T) {
	// get the issuing and distribution keypairs
	pairIss, pairDis, err := getAssetKeypairs()
	if err != nil {
		t.Error(err)
	}
	if err = colon.MAllowTrust(pairIss, "VEF", pairDis.Address(), true, false); err != nil {
		t.Error(err)
	}
}

func TestTransAllowTrust2(t *testing.T) {
	// get the issuing and distribution keypairs
	pairIss, pairDis, err := getAssetKeypairs()
	if err != nil {
		t.Error(err)
	}
	// create a transaction
	if tb, err1 := colon.MTrans(pairIss.Address()); err1 == nil {
		// add operation to the transaction
		opMut := build.AllowTrust(build.Trustor{pairDis.Address()}, build.AllowTrustAsset{Code: "VEF"}, build.Authorize{Value: false})
		if err = colon.MOpsAdd(tb, opMut); err == nil {
			// sign the transaction
			if txe, err2 := colon.MSign(tb, pairIss.Seed()); err2 == nil {
				// submit the transaction
				fmt.Println("AllowTrust Transaction", "VEF", "from", pairIss.Address(), "to", pairDis.Address(), "baseFee", tb.BaseFee)
				if resp, err3 := colon.MSubmit(txe); err3 == nil {
					fmt.Println("..successful", "Ledger", resp.Ledger, "Hash", resp.Hash)
				} else {
					txCode, opCodes, err := colon.MHorizonErrorResultCode(err3)
					fmt.Println("..failure", "txCode", txCode, "opCodes", opCodes, "err", err)
					err = err3
				}
			} else {
				err = err2
			}
		}
	} else {
		err = err1
	}
	if err != nil {
		t.Error(err)
	}
}

func TestTransAllowTrust3(t *testing.T) {
	// get the issuing and distribution keypairs
	pairIss, pairDis, err := getAssetKeypairs()
	if err != nil {
		t.Error(err)
	}
	// create a transaction
	opMut := build.AllowTrust(build.Trustor{pairDis.Address()}, build.AllowTrustAsset{Code: "VEF"}, build.Authorize{Value: false})
	if tb, err1 := colon.MTrans(pairIss.Address(), opMut); err1 == nil {
		// sign the transaction
		if txe, err2 := colon.MSign(tb, pairIss.Seed()); err2 == nil {
			// submit the transaction
			fmt.Println("AllowTrust Transaction", "VEF", "from", pairIss.Address(), "to", pairDis.Address(), "baseFee", tb.BaseFee)
			if resp, err3 := colon.MSubmit(txe); err3 == nil {
				fmt.Println("..successful", "Ledger", resp.Ledger, "Hash", resp.Hash)
			} else {
				txCode, opCodes, err := colon.MHorizonErrorResultCode(err3)
				fmt.Println("..failure", "txCode", txCode, "opCodes", opCodes, "err", err)
				err = err3
			}
		} else {
			err = err2
		}
	} else {
		err = err1
	}
	if err != nil {
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

func TestTransUnmarshal(t *testing.T) {
	var data string
	data = "AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAACgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAO5rKAAAAAAAAAAABVvwF9wAAAEAKZ7IPj/46PuWU6ZOtyMosctNAkXRNX9WCAI5RnfRk+AyxDLoDZP/9l3NvsxQtWj9juQOuoBlFLnWu8intgxQA"
	//data = "AAAAANHHHV431oTqGT5+8aPP6ugtg8KqLngW5mqIt08fj4zwAAAAyACK5vEAAAAKAAAAAAAAAAAAAAACAAAAAAAAAAAAAAAAuQqTymdxQpRPo74QkRHLtmL0dz3gGmm3botyqivAd/MAAAAABfXhAAAAAAAAAAABAAAAALkKk8pncUKUT6O+EJERy7Zi9Hc94Bppt26LcqorwHfzAAAAAAAAAAA7msoAAAAAAAAAAAEfj4zwAAAAQNqh3JSkniabqPrVemgwC2lew+gn5tL8c3eykM2BozUgHkXNLPNSq9d3b2C2kQGAWqSX+L7OSBEeiGSVZcXysQU="
	tx, err := colon.MXdrToTrans(data)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Tx source acc %d %x, fee %d, sequence %d, ext %d\n", tx.Tx.SourceAccount.Type, tx.Tx.SourceAccount.Ed25519, tx.Tx.Fee, tx.Tx.SeqNum, tx.Tx.Ext)
	fmt.Printf("Tx timebounds %v\n", tx.Tx.TimeBounds)
	fmt.Printf("Tx memo %v\n", tx.Tx.Memo)
	for i, o := range tx.Tx.Operations {
		fmt.Printf("Tx op%d %v\n", i, o)
	}
	for i, o := range tx.Signatures {
		fmt.Printf("Tx sig%d: hint %x, signature %x\n", i, o.Hint, o.Signature)
	}
}
