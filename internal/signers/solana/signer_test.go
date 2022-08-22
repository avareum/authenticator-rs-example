package solana

import (
	"context"
	"testing"

	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/test-go/testify/require"
)

func NewTestSigner() *SolanaSigner {
	s := NewSolanaSigner(SolanaSignerOptions{
		RPC: "http://127.0.0.1:8899",
	})
	s.Init()
	return s
}

func Test_SolanaSigner(t *testing.T) {
	t.Run("should sign payload", func(t *testing.T) {
		suite := fixtures.NewTestSuite()
		signer := NewTestSigner()
		receiver := solana.NewWallet()
		originalTx := suite.Solana.NewTx(system.NewTransferInstruction(
			100000,
			suite.Solana.Fund.PublicKey(),
			receiver.PublicKey(),
		).Build())

		t.Run("should sign relay tx", func(t *testing.T) {
			signatures, err := signer.sign(context.TODO(), originalTx, suite.Solana.Fund.PrivateKey)
			require.Nil(t, err)
			require.Nil(t, originalTx.VerifySignatures())
			require.Equal(t, 1, len(signatures))
		})

	})
}
