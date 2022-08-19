package wallet

import (
	"context"
	"fmt"

	"github.com/avareum/avareum-hubble-signer/internal/app"
	"github.com/avareum/avareum-hubble-signer/internal/signers"
	"github.com/avareum/avareum-hubble-signer/internal/signers/ethereum"
	"github.com/avareum/avareum-hubble-signer/internal/signers/solana"
	"github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/avareum/avareum-hubble-signer/pkg/acl"
	"github.com/avareum/avareum-hubble-signer/pkg/secret_manager"
	solanalib "github.com/gagliardetto/solana-go"
)

type WalletHandler interface {
	NewWallet() (*NewWalletResponse, error)
	Execute(req *types.SignerRequest) (*ExecuteWalletResponse, error)
}

type NewWalletResponse struct {
	Wallet string `json:"wallet"`
}

type ExecuteWalletResponse struct {
	Signatures []string `json:"signatures"`
}

type FundWalletHandler struct {
	WalletHandler
	app *app.AppSigner
}

func NewFundWalletHandler() (WalletHandler, error) {
	sm, err := secret_manager.NewGCPSecretManager()
	if err != nil {
		return nil, err
	}
	acl, err := acl.NewServiceACL()
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return &FundWalletHandler{
		app: app,
	}, err
}

func (f *FundWalletHandler) NewWallet() (*NewWalletResponse, error) {
	sm, err := secret_manager.NewGCPSecretManager()
	if err != nil {
		wallet := solanalib.NewWallet()
		walletNamespace := fmt.Sprintf("%s%s", signers.WALLET_PREFIX, wallet.PublicKey().String())
		sm.Create(walletNamespace, wallet.PrivateKey)
		return &NewWalletResponse{
			Wallet: wallet.PublicKey().String(),
		}, err
	}
	return nil, err
}

func (f *FundWalletHandler) Execute(req *types.SignerRequest) (*ExecuteWalletResponse, error) {
	res, err := f.app.TrySign(context.TODO(), *req)
	if err != nil {
		return nil, err
	}
	return &ExecuteWalletResponse{
		Signatures: res.Signatures,
	}, nil
}
