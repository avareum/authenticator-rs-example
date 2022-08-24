package ethereum

import (
	"fmt"
	"testing"

	"github.com/avareum/avareum-hubble-signer/internal/signers/ethereum/types"
	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/test-go/testify/require"
)

func Test_EthereumSigner(t *testing.T) {
	t.Run("should sign payload", func(t *testing.T) {
		suite := fixtures.NewTestSuite()
		signer := NewEthereumSigner(EthereumSignerOptions{})
		sender := types.MustNewEthereumKey()
		receiver := types.MustNewEthereumKey()
		originalTx := suite.Ethereum.NewTransferTransaction(
			crypto.PubkeyToAddress(sender.PublicKey),
			crypto.PubkeyToAddress(receiver.PublicKey),
			1,
		)

		t.Run("should sign relay tx", func(t *testing.T) {
			signedTx, err := signer.sign(originalTx, sender)
			require.Nil(t, err)
			require.NotNil(t, signedTx)

			t.Run("should successfully broadcast signed tx", func(t *testing.T) {
				suite.Ethereum.FaucetTo(sender.PublicKey)
				suite.Ethereum.SendTransaction(signedTx)

				receipt, err := suite.Ethereum.TransactionReceipt(signedTx.Hash())
				require.Nil(t, err)
				require.NotNil(t, receipt)
				fmt.Println(receipt)
			})

		})

	})
}
