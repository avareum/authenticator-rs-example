package solana

import (
	"context"
	"encoding/base64"
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

func Test_SolanaDecoder(t *testing.T) {
	t.Run("should decode tx", func(t *testing.T) {
		suite := fixtures.NewTestSuite()
		signer := NewTestSigner()
		receiver := solana.NewWallet()
		originalTx := suite.Solana.NewTx(system.NewTransferInstruction(
			100000,
			suite.Solana.Fund.PublicKey(),
			receiver.PublicKey(),
		).Build())

		t.Run("should decode transfer tx", func(t *testing.T) {
			rawTx, err := base64.StdEncoding.DecodeString(originalTx.Message.ToBase64())
			require.Nil(t, err)

			decodedTx, err := signer.decode(context.TODO(), rawTx)
			require.Nil(t, err)

			t.Run("should contain transfer instruction", func(t *testing.T) {
				require.Equal(t, 1, len(decodedTx.Message.Instructions))
			})

			t.Run("should contain transfer accounts", func(t *testing.T) {
				require.Equal(t, 3, len(decodedTx.Message.AccountKeys))
			})

			t.Run("should contain system account (transfer SOL)", func(t *testing.T) {
				program, err := decodedTx.ResolveProgramIDIndex(decodedTx.Message.Instructions[0].ProgramIDIndex)
				require.Nil(t, err)
				require.Equal(t, "11111111111111111111111111111111", program.String())
			})

			t.Run("should sign relay tx", func(t *testing.T) {
				signatures, err := signer.sign(context.TODO(), decodedTx, suite.Solana.Fund.PrivateKey)
				require.Nil(t, err)
				require.Nil(t, decodedTx.VerifySignatures())
				require.Equal(t, 1, len(signatures))
			})
		})

	})
}
