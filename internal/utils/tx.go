package utils

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func WaitSolanaTxConfirmed(client *rpc.Client, sig solana.Signature) {
	for {
		txMetadata, _ := client.GetConfirmedTransaction(context.TODO(), sig)
		if txMetadata != nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
		continue
	}

}
