package ethereum

import (
	"fmt"
	"testing"

	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/test-go/testify/require"
)

func Test_EthereumSigner(t *testing.T) {
	t.Run("should sign payload", func(t *testing.T) {
		suite := fixtures.NewTestSuite()
		signer := NewEthereumSigner(EthereumSignerOptions{})
		sender := suite.Ethereum.MustNewWallet()
		receiver := suite.Ethereum.MustNewWallet()
		originalTx := suite.Ethereum.NewTransferTransaction(*sender, receiver.PublicKey, 1)

		fmt.Println("sender", crypto.PubkeyToAddress(sender.PublicKey))
		fmt.Println("receiver", crypto.PubkeyToAddress(receiver.PublicKey))

		t.Run("should sign relay tx", func(t *testing.T) {
			signedTx, err := signer.sign(originalTx, sender)
			require.Nil(t, err)
			require.NotNil(t, signedTx)

			t.Run("should successfully broadcast signed tx", func(t *testing.T) {
				suite.Ethereum.AirdropTo(sender.PublicKey)
				suite.Ethereum.SendTransaction(signedTx)

				receipt, err := suite.Ethereum.TransactionReceipt(signedTx.Hash())
				require.Nil(t, err)
				require.NotNil(t, receipt)
			})

		})

	})
}
