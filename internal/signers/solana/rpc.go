package solana

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/avareum/avareum-hubble-signer/internal/signers/signer"
	"github.com/avareum/avareum-hubble-signer/internal/signers/types"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

type SolanaSigner struct {
	signer.BaseSigner
	opt       SolanaSignerOptions
	rpcclient *rpc.Client
	wsclient  *ws.Client
}

type SolanaSignerOptions struct {
	RPC       string
	Websocket string
}

// Signer implementation checked against internal/signers/types/signer.go
var _ types.Signer = &SolanaSigner{}

func NewSolanaSigner(opt SolanaSignerOptions) *SolanaSigner {
	s := &SolanaSigner{
		opt: opt,
	}
	return s
}

func (s *SolanaSigner) ID() string {
	return "solana.mainnet-beta"
}

func (s *SolanaSigner) Init() error {
	// create a new rpc & websocket client (used for confirming transactions)
	rpcClient := rpc.New(s.opt.RPC)
	wsClient, err := ws.Connect(context.Background(), s.opt.Websocket)
	if err != nil {
		return err
	}
	s.rpcclient = rpcClient
	s.wsclient = wsClient
	return s.BaseSigner.Init()
}

func (s *SolanaSigner) parseSignerKey() (solana.PrivateKey, error) {
	raw, err := s.BaseSigner.FetchSignerRawKey()
	if err != nil {
		return nil, err
	}

	var values []byte
	err = json.Unmarshal(raw, &values)
	if err != nil {
		return nil, fmt.Errorf("SolanaSigner: decode keygen file: %w", err)
	}
	return solana.PrivateKey([]byte(values)), nil
}

func (s *SolanaSigner) decodeTx(data []byte) (*solana.Transaction, error) {
	recent, err := s.rpcclient.GetRecentBlockhash(context.TODO(), rpc.CommitmentConfirmed)
	if err != nil {
		return nil, err
	}
	decodedTx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(data))
	if err != nil {
		return nil, err
	}
	decodedTx.Message.RecentBlockhash = recent.Value.Blockhash
	// payer must be mark as signer while tx building

	return decodedTx, nil
}

func (s *SolanaSigner) sign(tx *solana.Transaction, account solana.PrivateKey) ([]solana.Signature, error) {
	return tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if account.PublicKey().Equals(key) {
			return &account
		}
		return nil
	})
}

func (s *SolanaSigner) SignAndBroadcast(req types.SignerRequest) ([]string, error) {
	// fetch key from secret manager
	signerAccount, err := s.parseSignerKey()
	if err != nil {
		return nil, err
	}

	// decode & validate tx from data
	tx, err := s.decodeTx(req.Payload)
	if err != nil {
		return nil, err
	}

	// sign tx using secret signer
	_, err = s.sign(tx, signerAccount)
	if err != nil {
		return nil, err
	}

	// broadcast tx
	signature, err := confirm.SendAndConfirmTransaction(context.TODO(), s.rpcclient, s.wsclient, tx)
	if err != nil {
		return nil, err
	}

	return []string{signature.String()}, err
}
