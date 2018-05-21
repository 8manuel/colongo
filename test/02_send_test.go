package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stellar/go/clients/horizon"
)

func TestSendStreamPayments(t *testing.T) {
	const address = "GCIKXGDA4G2JWYMBOU27THPA6MSBRRVX4YDIJOKXP6KIIGZVGUVGA52A"
	ctx := context.Background()
	lastCursor := ""

	fmt.Println("Waiting for a payment...")

	for {

		cursor := horizon.Cursor(lastCursor)
		err := horizon.DefaultTestNetClient.StreamPayments(ctx, address, &cursor, func(payment horizon.Payment) {
			fmt.Println("Payment type", payment.Type)
			fmt.Println("Payment Paging Token", payment.PagingToken)
			fmt.Println("Payment From", payment.From)
			fmt.Println("Payment To", payment.To)
			fmt.Println("Payment Asset Type", payment.AssetType)
			fmt.Println("Payment Asset Code", payment.AssetCode)
			fmt.Println("Payment Asset Issuer", payment.AssetIssuer)
			fmt.Println("Payment Amount", payment.Amount)
			fmt.Println("Payment Memo Type", payment.Memo.Type)
			fmt.Println("Payment Memo", payment.Memo.Value)
			fmt.Println("")
			lastCursor = payment.PagingToken
		})

		if err != nil {
			fmt.Printf("Test continues even that there is the following error: %s\n", err)
			//t.Error(err)
		}
	}
}
