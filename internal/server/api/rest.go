package api

import (
	"net/http"

	"github.com/avareum/avareum-hubble-signer/internal/server/wallet"
	"github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/gin-gonic/gin"
)

type RestAPI struct{}

func NewRestAPI() *RestAPI {
	return &RestAPI{}
}

// Serve starts the simplest rest api server
func (api *RestAPI) Serve() {
	walletHandler, err := wallet.NewFundWalletHandler()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.POST("/wallet/new", func(c *gin.Context) {
			res, err := walletHandler.NewWallet()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status":   "OK",
					"response": res,
				})
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "FAILED",
				"error":  "create new wallet failed",
			})
		})

		v1.POST("/wallet/execute", func(c *gin.Context) {
			// TODO: parse request body
			res, err := walletHandler.Execute(&types.SignerRequest{})
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status":   "OK",
					"response": res,
				})
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "FAILED",
				"error":  "execute wallet failed",
			})
		})
	}
	r.Run()
}
