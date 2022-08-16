package app

import (
	"context"
	"testing"
	"time"

	"github.com/avareum/avareum-hubble-signer/internal/signers/solana"
	signerTypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/avareum/avareum-hubble-signer/internal/types"
	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	solanalib "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/stretchr/testify/assert"
)

func Test_App(t *testing.T) {
	assert := assert.New(t)
	suite := fixtures.NewTestSuite()
	sm := fixtures.NewTestSecretManager()
	mq := fixtures.NewTestMessageQueue()
	acl := fixtures.NewTestACL()
	acl.CreateTestServiceKey("caller-service")
	reqHandler := make(chan types.SignerRequestedResponse)

	// [hack] prepare mock tx & signer
	receiver := solanalib.NewWallet()
	originalTx := suite.Solana.NewTx(system.NewTransferInstruction(
		1*solanalib.LAMPORTS_PER_SOL,
		suite.Solana.Fund.PublicKey(),
		receiver.PublicKey(),
	).Build())

	// [hack] store fund key on secret manager
	sm.Create(suite.Solana.Fund.PublicKey().String(), suite.Solana.Fund.PrivateKey)

	app := NewAppSigner()
	app.RegisterSecretManager(sm)
	app.RegisterACL(acl)
	app.RegisterSignerRequestedResponseHandler(reqHandler)

	err := app.AddSigners(solana.NewSolanaSigner(solana.SolanaSignerOptions{
		RPC:       "http://127.0.0.1:8899",
		Websocket: "ws://localhost:8900",
	}))
	assert.Nil(err)

	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		// [hack] encode tx to binary
		payload, err := originalTx.Message.MarshalBinary()
		assert.Nil(err)

		// [hack] use service key sign payload
		payloadSignature, err := acl.SignPayload("caller-service", payload)
		assert.Nil(err)

		// [hack] push the request
		mq.Push(signerTypes.SignerRequest{
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

	// start long running receiving, signing, and broadcasting
	go app.Receive(ctx, mq)

	// wait for the response
	response := <-reqHandler
	assert.Nil(response.Error)
	assert.NotNil(response.Signatures)
	assert.Equal(1, len(*response.Signatures))
}
