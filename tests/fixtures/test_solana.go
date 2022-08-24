package fixtures

import (
	"context"
	"log"

	"github.com/avareum/avareum-hubble-signer/internal/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

const lamports = 1000000000

type SolanaTestSuite struct {
	Fund   *solana.Wallet
	client *rpc.Client
}

func NewSolanaTestSuite() *SolanaTestSuite {
	m := &SolanaTestSuite{
		Fund:   solana.NewWallet(),
		client: rpc.New("http://127.0.0.1:8899"),
	}
	return m
}

func (m *SolanaTestSuite) Sign(payload []byte) (solana.Signature, error) {
	return m.Fund.PrivateKey.Sign(payload)
}

func (m *SolanaTestSuite) Faucet() {
	m.FaucetTo(m.Fund.PublicKey())
}

func (m *SolanaTestSuite) FaucetTo(to solana.PublicKey) {
	sig, err := m.client.RequestAirdrop(context.TODO(), to, 10*lamports, rpc.CommitmentFinalized)
	if err != nil {
		log.Fatal(err)
	}
	utils.WaitSolanaTxConfirmed(m.client, sig)
}

func (m *SolanaTestSuite) NewTx(payer solana.PublicKey, ixs ...solana.Instruction) *solana.Transaction {
	recent, err := m.client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		log.Fatal(err)
	}
	tx, _ := solana.NewTransaction(
		ixs,
		recent.Value.Blockhash,
		solana.TransactionPayer(payer),
	)
	return tx
}
