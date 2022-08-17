package solana

import (
	"context"
	"fmt"

	"github.com/avareum/avareum-hubble-signer/internal/signers"
	"github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/avareum/avareum-hubble-signer/internal/utils"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type SolanaSigner struct {
	signers.BaseSigner
	opt       SolanaSignerOptions
	rpcclient *rpc.Client
}

type SolanaSignerOptions struct {
	RPC string
}

// Signer implementation checked against internal/signers/types/signer.go
var _ types.Signer = (*SolanaSigner)(nil)

func NewSolanaSigner(opt SolanaSignerOptions) *SolanaSigner {
	s := &SolanaSigner{
		opt: opt,
	}
	return s
}

// ID returns the signer's ID
func (s *SolanaSigner) ID() string {
	return "solana.mainnet-beta"
}

// Init create a new rpc & websocket client (used for confirming transactions)
func (s *SolanaSigner) Init() error {
	rpcClient := rpc.New(s.opt.RPC)
	s.rpcclient = rpcClient
	return s.BaseSigner.Init()
}

// SignTransaction sign a transaction with the signer's private key
func (s *SolanaSigner) SignAndBroadcast(ctx context.Context, req types.SignerRequest) ([]string, error) {
	fundSigner, err := s.getFundSignerKey(ctx, req.Fund)
	if err != nil {
		return nil, err
	}
	tx, err := s.decode(ctx, req.Payload)
	if err != nil {
		return nil, err
	}
	recent, err := s.rpcclient.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	tx.Message.RecentBlockhash = recent.Value.Blockhash
	_, err = s.sign(ctx, tx, fundSigner)
	if err != nil {
		return nil, err
	}
	signature, err := s.rpcclient.SendTransaction(ctx, tx)
	utils.WaitSolanaTxConfirmed(s.rpcclient, signature)
	return []string{signature.String()}, err
}

/*
 Internal
*/

func (s *SolanaSigner) getFundSignerKey(ctx context.Context, fund string) (solana.PrivateKey, error) {
	raw, err := s.BaseSigner.FetchSignerRawKey(fund)
	if err != nil {
		return nil, err
	}
	return solana.PrivateKey([]byte(raw)), nil
}

func (s *SolanaSigner) decode(ctx context.Context, payload []byte) (*solana.Transaction, error) {
	message := solana.Message{}
	err := bin.UnmarshalBin(&message, payload)
	if err != nil {
		return nil, fmt.Errorf("SolanaSigner: unmarshal tx msg failed: %v", err)
	}
	tx := solana.Transaction{}
	tx.Message = message
	// TODO: added feature partial sign
	return &tx, nil
}

func (s *SolanaSigner) sign(ctx context.Context, tx *solana.Transaction, account solana.PrivateKey) ([]solana.Signature, error) {
	signatures, err := tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if account.PublicKey().Equals(key) {
			return &account
		}
		return nil
	})
	// payer must be mark as signer while tx building
	return signatures, err
}
