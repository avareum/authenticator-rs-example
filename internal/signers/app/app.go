package app

import (
	"fmt"
	"os"

	"github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/avareum/avareum-hubble-signer/pkg/whitelist"
)

type App struct {
	Signers map[string]types.Signer
	wl      *whitelist.Whitelist
}

func NewApp() *App {
	a := &App{}
	return a
}

func (a *App) WithSigner(signers ...types.Signer) error {
	for _, s := range signers {
		err := s.Init()
		if err != nil {
			return err
		}
		a.Signers[s.ID()] = s
	}
	return nil
}

func (a *App) Start() error {
	// create whitelist checker
	wl, err := whitelist.NewWhitelist(whitelist.WhitelistOptions{
		ProjectID: os.Getenv("GCP_PROJECT"),
		Bucket:    "service-whitelists",
	})
	if err != nil {
		return err
	}
	a.wl = wl

	// TODO: initiate message queue connection
	receiver := make(chan types.SignerRequest)

	// Mock request
	go func() {
		// TODO: parse request from message queue to SignerRequest
		receiver <- types.NewMockSignerRequest()
	}()

	for {
		req := <-receiver

		// check if the caller is whitelisted
		if !a.wl.CanCall(req.Caller, req.Payload, req.Signature) {
			fmt.Println("Caller is not whitelisted")
			continue
		}

		if signer, isExists := a.Signers[req.SignerID()]; isExists {
			_, err := signer.SignAndBroadcast(req)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Signer not found")
		}

		// TODO: publish broadcasted signatures
		// TODO: ack mq
	}
}
