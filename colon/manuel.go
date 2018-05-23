package colon

import (
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/go-errors/errors"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

// baseSeed is the base seed used to generate deterministic keypairs, the lenght is 24bytes
var baseSeed string

// Init sets the package base seed used to generate deterministic keypairs
func Init(seed string) (err error) {
	l := len(seed)
	if l < 8 {
		return errors.New("Minimum seed length 8")
	} else if l < 24 {
		seed = seed + strings.Repeat(" ", 24-l)
		//seed = seed + ("                ")[0:24-l]
	} else if l > 24 {
		return errors.New("Maximum seed length 24")
	}
	baseSeed = seed
	return nil
}

// DeterministicKeypair generates a keypair (seed+address) using baseSeed+accName as part of the seed; accName maximum length is 8.
func DeterministicKeypair(accName string) (pair *keypair.Full) {
	// if accName length is longer than 8 return nil because this is an error
	if len(accName) > 8 || accName == "" {
		return nil
	}
	// generate a [32]byte seed for the issuing account
	bytes := []byte(baseSeed + accName)
	byteSeed := [32]byte{}
	bl := len(bytes)
	for i := 0; i < bl; i++ {
		byteSeed[i] = bytes[i]
	}
	// generate the keypair from the seed
	var err error
	if pair, err = keypair.FromRawSeed(byteSeed); err != nil {
		return nil
	}
	return pair
}

// MSeed2Bytes converts a seed into a [32]byte array
func MSeed2Bytes(seed string) (seedBytes [32]byte, err error) {
	// convert the seed to []byte
	data, err := base32.StdEncoding.DecodeString(seed)
	if err != nil {
		return seedBytes, err
	}
	// copy the bytes from 1 to 32
	for i := 0; i < 32; i++ {
		seedBytes[i] = data[i+1]
	}

	return seedBytes, err
}

// MLoadAccount gets the account data from the horizon server
func MLoadAccount(addr string) (account horizon.Account, err error) {
	if account, err = horizon.DefaultTestNetClient.LoadAccount(addr); err != nil {
		MHorizonProblemView(err)
	}
	return account, err
}

// MSign signs the transaction with all the signers provided.
func MSign(tx *build.TransactionBuilder, signers ...string) (txe build.TransactionEnvelopeBuilder, err error) {
	// It just calls the transaction sign function with all the signers provided.
	return tx.Sign(signers...)
}

// MSignAdd adds new signature to the transaction envelope; as is passed the envelope addres it does not return the envelope to the caller.
func MSignAdd(txe *build.TransactionEnvelopeBuilder, signer string) (err error) {
	// It just calls the transaction envelope builder mutate with a sign object constructed with the signer string.
	return txe.Mutate(build.Sign{signer})
}

// MSignSubmit signs the transaction, converts to base64 and sends to Stellar through horizon server.
func MSignSubmit(seed string, tx *build.TransactionBuilder) (resp horizon.TransactionSuccess, err error) {
	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(seed)
	if err != nil {
		return resp, err
	}
	// Convert to base64
	txeB64, err := txe.Base64()
	if err != nil {
		return resp, err
	}
	// Send to Stellar
	resp, err = horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		MHorizonProblemView(err)
		return resp, err
	}
	return resp, err
}

// MSubmit converts a transaction envelope builder to base64 and sends to Stellar through horizon server.
func MSubmit(seed string, txe build.TransactionEnvelopeBuilder) (resp horizon.TransactionSuccess, err error) {
	// Convert to base64
	txeB64, err := txe.Base64()
	if err != nil {
		return resp, err
	}
	// Send to Stellar
	resp, err = horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		MHorizonProblemView(err)
		return resp, err
	}
	return resp, err
}

// MHorizonProblemView prints the details of an horizon error to be able to view what happens
func MHorizonProblemView(err error) {
	eo := err.(*horizon.Error)
	eop := eo.Problem
	fmt.Println("Horizon Problem", "status", eop.Status)
	fmt.Println("- type:", eop.Type)
	fmt.Println("- title:", eop.Title)
	fmt.Println("- detail:", eop.Detail)
	fmt.Println("- instance:", eop.Instance)
	for k, v := range eop.Extras {
		fmt.Println("- map", k, ":", string(v))
	}
}

// MHorizonErrorResultCode extracts and returns from an error (that can be casted to horizon.Error) the transaction code and the operation codes.
func MHorizonErrorResultCode(herr error) (txCode string, opCodes []string, err error) {
	eo := herr.(*horizon.Error)
	if rc, err := eo.ResultCodes(); err != nil {
		return txCode, opCodes, err
	} else {
		return rc.TransactionCode, rc.OperationCodes, nil
	}
}

// MSetOptions sets the address options.
// The options to set are defined in opts map with the option name as key and the value as option value that has a different type depending on the option.
//  - InflationDest: is a [32]byte with the address publickey
//  - ClearFlags/SetFlags/MasterWeight/LowThreshold/MedThreshold/HighThreshold/HomeDOmain: is a uint32
//  - Signer: is an interface array with [keyType int32, address/transaction/hash int32, weight uint32)
func MSetOptions(pair *keypair.Full, opts map[string]interface{}) (err error) {
	// create and fill the SetOptions with the opts map (other way is to create muts:=[]interface{}, compose options and then build.SetOptions(muts...) to create it)
	so := build.SetOptions()
	for k, v := range opts {
		switch k {
		case "InflationDest":
			xpk, err := xdr.NewPublicKey(xdr.PublicKeyTypePublicKeyTypeEd25519, v.([32]byte))
			if err != nil {
				return err
			}
			x := xdr.AccountId(xpk)
			so.SO.InflationDest = &x
		case "ClearFlags":
			x := xdr.Uint32(v.(uint32))
			so.SO.ClearFlags = &x
		case "SetFlags":
			x := xdr.Uint32(v.(uint32))
			so.SO.SetFlags = &x
		case "MasterWeight":
			x := xdr.Uint32(v.(uint32))
			so.SO.MasterWeight = &x
		case "LowThreshold":
			x := xdr.Uint32(v.(uint32))
			so.SO.LowThreshold = &x
		case "MedThreshold":
			x := xdr.Uint32(v.(uint32))
			so.SO.MedThreshold = &x
		case "HighThreshold":
			x := xdr.Uint32(v.(uint32))
			so.SO.HighThreshold = &x
		case "HomeDomain":
			x := xdr.String32(v.(string))
			so.SO.HomeDomain = &x
		case "Signer":
			xa := v.([]interface{})
			xs, err := xdr.NewSignerKey(xdr.SignerKeyType(xa[0].(int32)), xdr.Uint256(xa[1].([32]byte)))
			if err != nil {
				return err
			}
			x := xdr.Signer{Key: xs, Weight: xdr.Uint32(xa[2].(uint32))}
			so.SO.Signer = &x
		default:
			return errors.New("wrong option")
		}
	}

	// compose the setOptions trust transaction
	seedDis := pair.Seed()
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{pair.Address()},
		build.AutoSequence{horizon.DefaultTestNetClient},
		so,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Sign and submit the transaction
	fmt.Println("SerOptions Transaction", "addr", pair.Address())
	if resp, err := MSignSubmit(seedDis, tx); err == nil {
		fmt.Println("..successful", "Ledger", resp.Ledger, "Hash", resp.Hash)
	}
	return err
}

// MTransPayment sends a payment transaction of amtStr from a pairSource address to a destination address.
// Instead of using directly the source seed it is used the pairSource, in this way the seed is used for signing and the address for displaying.
// If checkDest is set then the destination account is verified before sending (so no fee is paid if the address not exists).
func MTransPayment(pairSource *keypair.Full, addrDest, asset, amtStr string, checkDest bool) (err error) {
	// Make sure destination address exists, so no fees are paid if it does not exist
	if checkDest {
		if _, err := horizon.DefaultTestNetClient.LoadAccount(addrDest); err != nil {
			panic(err)
		}
	}

	// Build the transaction
	var pb build.PaymentBuilder
	if asset == "" {
		pb = build.Payment(build.Destination{addrDest}, build.NativeAmount{amtStr})
	} else {
		pb = build.Payment(build.Destination{addrDest}, build.CreditAmount{asset, pairSource.Address(), amtStr})
	}
	seedSource := pairSource.Seed()
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{seedSource},
		build.AutoSequence{horizon.DefaultTestNetClient},
		pb,
	)
	if err != nil {
		return err
	}
	// Sign and submit the transaction
	fmt.Println("Payment Transaction", asset, amtStr, "from", pairSource.Address(), "to", addrDest)
	if resp, err := MSignSubmit(seedSource, tx); err == nil {
		fmt.Println("..successful", "Ledger", resp.Ledger, "Hash", resp.Hash)
	}
	return err
}

func MTransTrust(pairDis *keypair.Full, assCode, addrIss, limitStr string, checkIss bool) (err error) {
	// Make sure issuing address (addrIss) exists, so no fees are paid if it does not exist
	if checkIss {
		if _, err := horizon.DefaultTestNetClient.LoadAccount(addrIss); err != nil {
			panic(err)
		}
	}

	// compose the trust transaction
	seedDis := pairDis.Seed()
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{seedDis},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.Trust(assCode, addrIss, build.Limit(limitStr)),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Sign and submit the transaction
	fmt.Println("Trust Transaction", assCode, limitStr, "from", pairDis.Address(), "to", addrIss)
	if resp, err := MSignSubmit(seedDis, tx); err == nil {
		fmt.Println("..successful", "Ledger", resp.Ledger, "Hash", resp.Hash)
	}
	return err
}

func MAllowTrust(pairIss *keypair.Full, assCode, addr string, authorize, checkAddr bool) (err error) {
	// Make sure address exists, so no fees are paid if it does not exist
	if checkAddr {
		if _, err := horizon.DefaultTestNetClient.LoadAccount(addr); err != nil {
			panic(err)
		}
	}

	// compose the allow trust transaction
	seedDis := pairIss.Seed()
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{pairIss.Address()},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.AllowTrust(build.Trustor{addr}, build.AllowTrustAsset{Code: assCode}, build.Authorize{Value: authorize}),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Sign and submit the transaction
	fmt.Println("AllowTrust Transaction", assCode, "from", pairIss.Address(), "to", addr)
	if resp, err := MSignSubmit(seedDis, tx); err == nil {
		fmt.Println("..successful", "Ledger", resp.Ledger, "Hash", resp.Hash)
	}
	return err
}
