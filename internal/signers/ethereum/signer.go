package ethereum

import (
	"context"
	"crypto/ecdsa"

	"github.com/avareum/avareum-hubble-signer/internal/signers"
	"github.com/avareum/avareum-hubble-signer/internal/signers/types"
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
var _ types.Signer = (*EthereumSigner)(nil)

func NewEthereumSigner(opt EthereumSignerOptions) *EthereumSigner {
	s := &EthereumSigner{
		opt:     opt,
		decoder: NewEthereumTransactionDecoder(),
	}
	return s
}

func (s *EthereumSigner) ID() string {
	return "ethereum.1"
}

func (s *EthereumSigner) Init() error {
	client, err := ethclient.Dial(s.opt.RPC)
	if err != nil {
		return err
	}
	s.client = client
	return nil
}

func (s *EthereumSigner) SignAndBroadcast(ctx context.Context, req types.SignerRequest) ([]string, error) {
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
