package solana

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/stretchr/testify/assert"
)

func NewTestSigner() *SolanaSigner {
	s := NewSolanaSigner(SolanaSignerOptions{
		RPC:       "http://127.0.0.1:8899",
		Websocket: "ws://localhost:8900",
	})
	s.Init()
	return s
}

var suite = fixtures.NewTestSuite()
var signer = NewTestSigner()

func Test_SignerDecoder(t *testing.T) {
	assert := assert.New(t)

	t.Run("should decode tx", func(t *testing.T) {
		receiver := solana.NewWallet()
		originalTx := suite.Solana.NewTx(system.NewTransferInstruction(
			100000,
			suite.Solana.Fund.PublicKey(),
			receiver.PublicKey(),
		).Build())

		t.Run("should decode transfer tx", func(t *testing.T) {
			rawTx, err := base64.StdEncoding.DecodeString(originalTx.Message.ToBase64())
			assert.Nil(err)

			decodedTx, err := signer.decode(context.TODO(), rawTx)
			assert.Nil(err)

			t.Run("should contain transfer instruction", func(t *testing.T) {
				assert.Equal(1, len(decodedTx.Message.Instructions))
			})

			t.Run("should contain transfer accounts", func(t *testing.T) {
				assert.Equal(3, len(decodedTx.Message.AccountKeys))
			})

			t.Run("should contain system account (transfer SOL)", func(t *testing.T) {
				program, err := decodedTx.ResolveProgramIDIndex(decodedTx.Message.Instructions[0].ProgramIDIndex)
				assert.Nil(err)
				assert.Equal("11111111111111111111111111111111", program.String())
			})

			t.Run("should sign relay tx", func(t *testing.T) {
				signatures, err := signer.sign(context.TODO(), decodedTx, suite.Solana.Fund.PrivateKey)
				assert.Nil(err)
				assert.Nil(decodedTx.VerifySignatures())
				assert.Equal(1, len(signatures))
			})
		})

	})
}
