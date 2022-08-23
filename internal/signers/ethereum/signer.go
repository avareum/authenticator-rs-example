package ethereum

import (
	"context"
	"crypto/ecdsa"

	"github.com/avareum/avareum-hubble-signer/constant"
	"github.com/avareum/avareum-hubble-signer/internal/signers"
	signertypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/avareum/avareum-hubble-signer/internal/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumSigner struct {
	signers.BaseSigner
	opt     EthereumSignerOptions
	decoder *EthereumTransactionDecoder
	client  *ethclient.Client
}

type EthereumSignerOptions struct {
	RPC string
}

// Signer implementation checked against internal/signers/types/signer.go
var _ signertypes.Signer = (*EthereumSigner)(nil)

func NewEthereumSigner(opt EthereumSignerOptions) *EthereumSigner {
	s := &EthereumSigner{
		opt:     opt,
		decoder: NewEthereumTransactionDecoder(),
	}
	return s
}

// Chain returns the signer's chain
func (s *EthereumSigner) Chain() types.Chain {
	return constant.EthereumMainnet
}

// Init create a new rpc client
func (s *EthereumSigner) Init() error {
	client, err := ethclient.Dial(s.opt.RPC)
	if err != nil {
		return err
	}
	s.client = client
	return nil
}

// SignTransaction sign a transaction with the signer's private key
func (s *EthereumSigner) SignAndBroadcast(ctx context.Context, req signertypes.SignerRequest) ([]string, error) {
	priv, err := s.getSigningKey(ctx, req.Wallet)
	if err != nil {
		return nil, err
	}
	tx, err := s.decoder.TryDecode(req.Payload)
	if err != nil {
		return nil, err
	}
	signedTx, err := s.sign(tx, priv)
	if err != nil {
		return nil, err
	}
	err = s.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	return []string{signedTx.Hash().Hex()}, nil
}

/* Internal */

func (s *EthereumSigner) getSigningKey(ctx context.Context, wallet string) (*ecdsa.PrivateKey, error) {
	raw, err := s.BaseSigner.FetchSignerRawKey(wallet)
	if err != nil {
		return nil, err
	}
	priv, err := crypto.ToECDSA(raw)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func (s *EthereumSigner) sign(tx *ethtypes.Transaction, priv *ecdsa.PrivateKey) (*ethtypes.Transaction, error) {
	return ethtypes.SignTx(tx, ethtypes.HomesteadSigner{}, priv)
}
