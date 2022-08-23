package main

import (
	"os"

	"github.com/avareum/avareum-hubble-signer/internal/server/api"
	"github.com/avareum/avareum-hubble-signer/pkg/logger"
)

func main() {
	// Override the default logger with a GCP logger.
	if os.Getenv("LOCAL_LOGGER") != "true" {
		gcpLogger, err := logger.NewGCPCloudLogger("avareum-hubble-signer")
		if err != nil {
			panic(err)
		}
		logger.Default = gcpLogger
	}

	// Create the app signer.
	api.NewRestAPI().Serve()
}
