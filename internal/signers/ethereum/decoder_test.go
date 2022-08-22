package ethereum

import (
	"testing"

	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	"github.com/test-go/testify/require"
)

func Test_EthereumTransactionDecoder(t *testing.T) {
	t.Run("ethereum decoder", func(t *testing.T) {
		suite := fixtures.NewTestSuite()
		decoder := NewEthereumTransactionDecoder()
		sender := suite.Ethereum.MustNewWallet()
		receiver := suite.Ethereum.MustNewWallet()
		originalTx := suite.Ethereum.NewTransferTransaction(*sender, receiver.PublicKey, 1)

		t.Run("should decode transaction", func(t *testing.T) {
			bin, err := originalTx.MarshalBinary()
			require.Nil(t, err)

			tx, err := decoder.DecodeFromTransaction(bin)
			require.Nil(t, err)
			require.NotNil(t, tx)
			require.Equal(t, originalTx.Hash(), tx.Hash())
		})

	})
}
