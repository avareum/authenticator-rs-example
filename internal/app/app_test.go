package app

import (
	"context"
	"testing"

	"github.com/avareum/avareum-hubble-signer/internal/signers/solana"
	signerTypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/avareum/avareum-hubble-signer/internal/types"
	smtypes "github.com/avareum/avareum-hubble-signer/pkg/secret_manager/types"
	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	solanalib "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/test-go/testify/require"
)

func NewTestTxPayload(suite *fixtures.TestSuite) []byte {
	receiver := solanalib.NewWallet()
	fund := suite.Solana.Fund.PublicKey()
	originalTx := suite.Solana.NewTx(fund, system.NewTransferInstruction(
		1*solanalib.LAMPORTS_PER_SOL,
		fund,
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
		_, err := app.TrySign(context.TODO(), suite.NewTestSignerRequest())
		require.EqualError(t, err, "secret manager is not registered")
	})

	t.Run("should reject invalid signer id requested", func(t *testing.T) {
		suite := fixtures.NewTestSuite()

		app := NewAppSigner().WithSecretManager(suite.SecretManager)
		_, err := app.TrySign(context.TODO(), signerTypes.SignerRequest{
			Chain:     types.NewChain("solono", "mainnet-beta"),
			Caller:    "",
			Wallet:    "",
			Payload:   []byte{},
			Signature: []byte{},
		})

		require.Regexp(t, "signer .* not found", err.Error())
	})

	t.Run("should reject mismatch service caller", func(t *testing.T) {
		suite := fixtures.NewTestSuite()
		suite.ACL.CreateTestServiceKey("caller-service")
		suite.ACL.CreateTestServiceKey("unauthorize-service")

		// [hack] use mismatch service key sign payload
		payload := NewTestTxPayload(suite)
		mismatchSignature, err := suite.ACL.SignPayload("unauthorize-service", payload)
		require.Nil(t, err)

		app := NewAppSigner().WithSecretManager(suite.SecretManager).WithACL(suite.ACL)
		_, err = app.TrySign(context.TODO(), signerTypes.SignerRequest{
			Chain:     types.NewChain("solana", "mainnet-beta"),
			Caller:    "caller-service",
			Wallet:    suite.Solana.Fund.PublicKey().String(),
			Payload:   payload,
			Signature: mismatchSignature,
		})
		require.EqualError(t, err, "invalid caller signature")
	})

	t.Run("should sign & broadcast valid request", func(t *testing.T) {
		suite := fixtures.NewTestSuite()
		suite.Solana.Faucet()

		// [hack] store fund key on secret manager
		suite.SecretManager.Create(
			smtypes.NewSecretWallet(suite.Solana.Fund.PublicKey().String()),
			suite.Solana.Fund.PrivateKey,
		)

		// [hack] create service key on ACL
		suite.ACL.CreateTestServiceKey("caller-service")

		// [hack] use valid service key sign payload
		payload := NewTestTxPayload(suite)
		payloadSignature, err := suite.ACL.SignPayload("caller-service", payload)
		require.Nil(t, err)

		app := NewAppSigner().WithSecretManager(suite.SecretManager).WithACL(suite.ACL)
		err = app.AddSigners(solana.NewSolanaSigner(solana.SolanaSignerOptions{
			RPC:   "http://127.0.0.1:8899",
			Chain: types.NewChain("solana", "localnet"),
		}))
		require.Nil(t, err)

		// start long running receiving, signing, and broadcasting
		response, err := app.TrySign(context.TODO(), signerTypes.SignerRequest{
			Chain:     types.NewChain("solana", "localnet"),
			Caller:    "caller-service",
			Wallet:    suite.Solana.Fund.PublicKey().String(),
			Payload:   payload,
			Signature: payloadSignature,
		})

		// wait for the response
		require.Nil(t, err)
		require.NotNil(t, response.Signatures)
		require.Equal(t, 1, len(response.Signatures))
	})
}
