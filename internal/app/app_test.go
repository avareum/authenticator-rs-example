package app

import (
	"context"
	"testing"
	"time"

	"github.com/avareum/avareum-hubble-signer/internal/signers/solana"
	signerTypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	solanalib "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/test-go/testify/require"
)

func NewTestTxPayload(suite *fixtures.TestSuite) []byte {
	receiver := solanalib.NewWallet()
	originalTx := suite.Solana.NewTx(system.NewTransferInstruction(
		1*solanalib.LAMPORTS_PER_SOL,
		suite.Solana.Fund.PublicKey(),
		receiver.PublicKey(),
	).Build())
	payload, err := originalTx.Message.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return payload
}

func Test_App(t *testing.T) {

	t.Run("should panic when missing secret manager", func(t *testing.T) {
		suite := fixtures.NewTestSuite()
		app := NewAppSigner()
		require.Panics(t, func() {
			app.Receive(context.TODO(), suite.MessageQueue)
		}, "should panic")
	})

	t.Run("should reject invalid request signer id", func(t *testing.T) {
		suite := fixtures.NewTestSuite()
		reqHandler := suite.NewSignerRequestedResponse()

		app := NewAppSigner()
		app.RegisterSecretManager(suite.SecretManager)
		app.RegisterSignerRequestedResponseHandler(reqHandler)
		go app.Receive(context.TODO(), suite.MessageQueue)

		// [hack] push mock request
		go suite.MessageQueue.Push(signerTypes.SignerRequest{
			Chain:     "solono",
			ChainID:   "mainnet-beta",
			Caller:    "",
			Fund:      "",
			Payload:   []byte{},
			Signature: []byte{},
		})

		response := <-reqHandler
		require.Regexp(t, "signer .* not found", response.Error.Error())
	})

	t.Run("should reject mismatch service caller", func(t *testing.T) {
		suite := fixtures.NewTestSuite()
		reqHandler := suite.NewSignerRequestedResponse()
		suite.ACL.CreateTestServiceKey("caller-service")
		suite.ACL.CreateTestServiceKey("unauthorize-service")

		// [hack] use mismatch service key sign payload
		payload := NewTestTxPayload(suite)
		mismatchSignature, err := suite.ACL.SignPayload("unauthorize-service", payload)
		require.Nil(t, err)

		app := NewAppSigner()
		app.RegisterSecretManager(suite.SecretManager)
		app.RegisterACL(suite.ACL)
		app.RegisterSignerRequestedResponseHandler(reqHandler)
		go app.Receive(context.TODO(), suite.MessageQueue)

		// [hack] push mock request
		go suite.MessageQueue.Push(signerTypes.SignerRequest{
			Chain:     "solana",
			ChainID:   "mainnet-beta",
			Caller:    "caller-service",
			Fund:      suite.Solana.Fund.PublicKey().String(),
			Payload:   payload,
			Signature: mismatchSignature,
		})

		response := <-reqHandler
		require.Error(t, response.Error, "invalid caller signature")
	})

	t.Run("should sign & broadcast valid request", func(t *testing.T) {
		suite := fixtures.NewTestSuite()

		// [hack] store fund key on secret manager
		suite.SecretManager.Create(suite.Solana.Fund.PublicKey().String(), suite.Solana.Fund.PrivateKey)

		// [hack] create service key on ACL
		suite.ACL.CreateTestServiceKey("caller-service")

		// [hack] use valid service key sign payload
		payload := NewTestTxPayload(suite)
		payloadSignature, err := suite.ACL.SignPayload("caller-service", payload)
		require.Nil(t, err)

		// [hack] prepare cancel context
		ctx, cancelCtx := context.WithCancel(context.Background())
		go func() {
			// [hack] push the request
			suite.MessageQueue.Push(signerTypes.SignerRequest{
				Chain:     "solana",
				ChainID:   "mainnet-beta",
				Caller:    "caller-service",
				Fund:      suite.Solana.Fund.PublicKey().String(),
				Payload:   payload,
				Signature: payloadSignature,
			})

			// wait for the response
			time.Sleep(5 * time.Second)
			cancelCtx()
		}()

		reqHandler := suite.NewSignerRequestedResponse()

		app := NewAppSigner()
		app.RegisterSecretManager(suite.SecretManager)
		app.RegisterACL(suite.ACL)
		app.RegisterSignerRequestedResponseHandler(reqHandler)
		err = app.AddSigners(solana.NewSolanaSigner(solana.SolanaSignerOptions{
			RPC: "http://127.0.0.1:8899",
		}))
		require.Nil(t, err)

		// start long running receiving, signing, and broadcasting
		go app.Receive(ctx, suite.MessageQueue)

		// wait for the response
		response := <-reqHandler
		require.Nil(t, response.Error)
		require.NotNil(t, response.Signatures)
		require.Equal(t, 1, len(*response.Signatures))
	})
}
