package solana

import (
	"context"
	"testing"

	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	"github.com/avareum/avareum-hubble-signer/tests/utils"
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

		t.Run("should decode encoded tx (partial signed)", func(t *testing.T) {
			bin, err := utils.MarshalBinarySolanaTransaction(originalTx)
			require.Nil(t, err)

			tx, err := signer.tryDecode(context.TODO(), bin)
			require.Nil(t, err)
			require.NotNil(t, tx)

			t.Run("should contain instruction", func(t *testing.T) {
				require.Equal(t, 1, len(tx.Message.Instructions))
			})

			t.Run("should contain accounts", func(t *testing.T) {
				require.Equal(t, 3, len(tx.Message.AccountKeys))
			})

			t.Run("should contain system account (transfer SOL)", func(t *testing.T) {
				program, err := tx.ResolveProgramIDIndex(tx.Message.Instructions[0].ProgramIDIndex)
				require.Nil(t, err)
				require.Equal(t, "11111111111111111111111111111111", program.String())
			})

			t.Run("should sign relay tx", func(t *testing.T) {
				signatures, err := signer.sign(context.TODO(), tx, suite.Solana.Fund.PrivateKey)
				require.Nil(t, err)
				require.Nil(t, tx.VerifySignatures())
				require.Equal(t, 1, len(signatures))
			})
		})

		t.Run("should decode encode message (only message data)", func(t *testing.T) {
			bin, err := originalTx.Message.MarshalBinary()
			require.Nil(t, err)

			tx, err := signer.tryDecode(context.TODO(), bin)
			require.Nil(t, err)
			require.NotNil(t, tx)

			t.Run("should contain instruction", func(t *testing.T) {
				require.Equal(t, 1, len(tx.Message.Instructions))
			})

			t.Run("should contain accounts", func(t *testing.T) {
				require.Equal(t, 3, len(tx.Message.AccountKeys))
			})

			t.Run("should contain system account (transfer SOL)", func(t *testing.T) {
				program, err := tx.ResolveProgramIDIndex(tx.Message.Instructions[0].ProgramIDIndex)
				require.Nil(t, err)
				require.Equal(t, "11111111111111111111111111111111", program.String())
			})

			t.Run("should sign relay tx", func(t *testing.T) {
				signatures, err := signer.sign(context.TODO(), tx, suite.Solana.Fund.PrivateKey)
				require.Nil(t, err)
				require.Nil(t, tx.VerifySignatures())
				require.Equal(t, 1, len(signatures))
			})
		})

	})
}
