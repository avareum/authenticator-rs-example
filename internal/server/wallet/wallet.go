package wallet

import (
	"context"
	"fmt"
	"os"

	"github.com/avareum/avareum-hubble-signer/constant"
	"github.com/avareum/avareum-hubble-signer/internal/app"
	"github.com/avareum/avareum-hubble-signer/internal/signers/ethereum"
	ethtypes "github.com/avareum/avareum-hubble-signer/internal/signers/ethereum/types"
	"github.com/avareum/avareum-hubble-signer/internal/signers/solana"
	signertypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/avareum/avareum-hubble-signer/internal/types"
	"github.com/avareum/avareum-hubble-signer/pkg/acl"
	"github.com/avareum/avareum-hubble-signer/pkg/secret_manager"
	smtypes "github.com/avareum/avareum-hubble-signer/pkg/secret_manager/types"
	"github.com/ethereum/go-ethereum/crypto"
	solanalib "github.com/gagliardetto/solana-go"
)

type WalletHandler interface {
	NewWallet(req NewWalletRequest) (*NewWalletResponse, error)
	Execute(req ExecuteWalletRequest) (*ExecuteWalletResponse, error)
}

type NewWalletRequest struct {
	Chain types.Chain
}

type NewWalletResponse struct {
	Wallet string `json:"wallet"`
}

type ExecuteWalletRequest struct {
	Chain          types.Chain
	SigningRequest signertypes.SignerRequest `json:"signing_request"`
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
			RPC: os.Getenv("ETHEREUM_ENDPOINT"),
		}),
		solana.NewSolanaSigner(solana.SolanaSignerOptions{
			RPC: os.Getenv("SOLANA_ENDPOINT"),
		}),
	)
	if err != nil {
		return nil, err
	}
	return &FundWalletHandler{
		app: app,
	}, err
}

func (f *FundWalletHandler) NewWallet(req NewWalletRequest) (*NewWalletResponse, error) {
	sm, err := secret_manager.NewGCPSecretManager()
	if err != nil {
		return nil, err
	}

	var priv []byte
	var wallet string
	// Create a new keypair of specific chain.
	// Derive keypair priv into raw key.
	// Store raw key as a payload in secret manager.
	// Label the secret with the prefix and wallet name which `WALLET_{wallet}`.
	switch req.Chain.ID() {
	case constant.EthereumMainnet.ID():
		ethKey, err := ethtypes.NewEthereumKey()
		if err != nil {
			return nil, err
		}
		priv = crypto.FromECDSA(ethKey)
		wallet = crypto.PubkeyToAddress(ethKey.PublicKey).Hex()
	case constant.SolanaMainnetBeta.ID():
		solanaKey, err := solanalib.NewRandomPrivateKey()
		if err != nil {
			return nil, err
		}
		priv = solanaKey
		wallet = solanaKey.PublicKey().String()
	default:
		return nil, fmt.Errorf("unknown chain %s", req.Chain.ID())
	}
	err = sm.Create(smtypes.NewSecretWallet(wallet), priv)
	if err != nil {
		return nil, err
	}
	return &NewWalletResponse{
		Wallet: wallet,
	}, nil
}

func (f *FundWalletHandler) Execute(req ExecuteWalletRequest) (*ExecuteWalletResponse, error) {
	res, err := f.app.TrySign(context.TODO(), req.SigningRequest)
	if err != nil {
		return nil, err
	}
	return &ExecuteWalletResponse{
		Signatures: res.Signatures,
	}, nil
}
