package solana

import (
	"testing"

	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	"github.com/avareum/avareum-hubble-signer/tests/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/test-go/testify/require"
)

func Test_SolanaTransactionDecoder(t *testing.T) {
	t.Run("solana decoder", func(t *testing.T) {
		suite := fixtures.NewTestSuite()
		decoder := NewSolanaTransactionDecoder()
		receiver := solana.NewWallet()
		originalTx := suite.Solana.NewTx(system.NewTransferInstruction(
			100000,
			suite.Solana.Fund.PublicKey(),
			receiver.PublicKey(),
		).Build())

		t.Run("should decode transaction", func(t *testing.T) {
			bin, err := utils.MarshalBinarySolanaTransaction(originalTx)
			require.Nil(t, err)

			tx, err := decoder.DecodeFromTransaction(bin)
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
		})

		t.Run("should try decode transaction", func(t *testing.T) {
			bin, err := utils.MarshalBinarySolanaTransaction(originalTx)
			require.Nil(t, err)

			tx, err := decoder.TryDecode(bin)
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
		})

		t.Run("should decode message bin", func(t *testing.T) {
			bin, err := originalTx.Message.MarshalBinary()
			require.Nil(t, err)

			tx, err := decoder.DecodeFromBinary(bin)
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
		})

		t.Run("should try decode message bin", func(t *testing.T) {
			bin, err := originalTx.Message.MarshalBinary()
			require.Nil(t, err)

			tx, err := decoder.TryDecode(bin)
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
		})

	})
}
