package main

import (
	"github.com/avareum/avareum-hubble-signer/internal/signers/app"
	"github.com/avareum/avareum-hubble-signer/internal/signers/ethereum"
	"github.com/avareum/avareum-hubble-signer/internal/signers/solana"
)

func main() {
	app := app.NewApp()
	err := app.WithSigner(
		ethereum.NewEthereumSigner(ethereum.EthereumSignerOptions{RPC: "https://rpc.ankr.com/eth"}),
		solana.NewSolanaSigner(solana.SolanaSignerOptions{
			RPC:       "https://api.devnet.solana.com",
			Websocket: "wss://api.devnet.solana.com",
		}),
	)
	if err != nil {
		panic(err)
	}
	app.Start()
}
