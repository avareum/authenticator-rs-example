package constant

import "github.com/avareum/avareum-hubble-signer/internal/types"

var (
	EthereumMainnet   = types.NewChain("ethereum", "1")
	SolanaMainnetBeta = types.NewChain("solana", "mainnet-beta")
	BSCTestnet        = types.NewChain("bsc", "97")
	SolanaDevnet      = types.NewChain("solana", "devnet")
)
