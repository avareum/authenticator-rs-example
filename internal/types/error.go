package types

import (
	signerTypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
)

type SignerRequestedResponse struct {
	Request    signerTypes.SignerRequest
	Signatures []string
}

type SignerRequestedResponseHandler = chan SignerRequestedResponse
