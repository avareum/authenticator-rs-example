package main

import (
	"github.com/avareum/avareum-hubble-signer/internal/app"
	"github.com/avareum/avareum-hubble-signer/internal/signers/ethereum"
	"github.com/avareum/avareum-hubble-signer/internal/signers/solana"
	"github.com/avareum/avareum-hubble-signer/pkg/acl"
	"github.com/avareum/avareum-hubble-signer/pkg/secret_manager"
)

func main() {
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

	// app.Receive()
}
