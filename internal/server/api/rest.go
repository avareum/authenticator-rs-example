package api

import (
	"fmt"
	"net/http"

	"github.com/avareum/avareum-hubble-signer/pkg/logger"
	"github.com/avareum/avareum-hubble-signer/pkg/secret_manager"
	"github.com/gagliardetto/solana-go"
	"github.com/gin-gonic/gin"
)

type RestAPI struct{}

func NewRestAPI() *RestAPI {
	return &RestAPI{}
}

// Serve starts the simplest rest api server
func (api *RestAPI) Serve() {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.POST("/wallet/create", func(c *gin.Context) {
			sm, err := secret_manager.NewGCPSecretManager()
			if err != nil {
				wallet := solana.NewWallet()
				walletNamespace := fmt.Sprintf("WALLET_%s", wallet.PublicKey().String())
				sm.Create(walletNamespace, wallet.PrivateKey)
				c.JSON(http.StatusOK, gin.H{
					"status": "OK",
					"wallet": wallet.PublicKey().String(),
				})
			}

			logger.Default.Err(fmt.Sprintf("API: /fund/wallet/new: %v", err))

			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "FAILED",
				"error":  "create new wallet failed",
			})
		})
	}
	r.Run()
}
