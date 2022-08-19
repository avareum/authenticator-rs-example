package main

import (
	"context"
	"os"

	"github.com/avareum/avareum-hubble-signer/internal/app"
	"github.com/avareum/avareum-hubble-signer/internal/message_queue"
	"github.com/avareum/avareum-hubble-signer/internal/signers/ethereum"
	"github.com/avareum/avareum-hubble-signer/internal/signers/solana"
	"github.com/avareum/avareum-hubble-signer/pkg/acl"
	"github.com/avareum/avareum-hubble-signer/pkg/logger"
	"github.com/avareum/avareum-hubble-signer/pkg/secret_manager"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	// Override the default logger with a GCP logger.
	gcpLogger, err := logger.NewGCPCloudLogger("avareum-hubble-signer")
	if err != nil {
		panic(err)
	}
	logger.Default = gcpLogger

	// Create the app signer.
	sm, err := secret_manager.NewGCPSecretManager()
	if err != nil {
		panic(err)
	}
	acl, err := acl.NewServiceACL()
	if err != nil {
		panic(err)
	}
	app := app.NewAppSigner()
	app.RegisterACL(acl)
	app.RegisterSecretManager(sm)
	err = app.AddSigners(
		ethereum.NewEthereumSigner(ethereum.EthereumSignerOptions{
			RPC: "https://rpc.ankr.com/eth",
		}),
		solana.NewSolanaSigner(solana.SolanaSignerOptions{
			RPC: "https://api.devnet.solana.com",
		}),
	)
	if err != nil {
		panic(err)
	}

	// Create the message queue.
	mq, err := message_queue.NewPubsubWithOpt(message_queue.PubsubOptions{
		SubscriptionID: os.Getenv("MQ_RECEIVE_CHANNEL"),
	})
	if err != nil {
		panic(err)
	}
	receiver := mq.ReceiveChannel()
	for {
		select {
		case req := <-receiver:
			response, err := app.TrySign(context.TODO(), req)
			if err != nil {
				logger.Default.Err(req, err)
			} else {
				logger.Default.Info(response)
			}
		}

	}
}
