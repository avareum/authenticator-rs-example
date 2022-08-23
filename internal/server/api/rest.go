package api

import (
	"net/http"

	"github.com/avareum/avareum-hubble-signer/internal/server/wallet"
	"github.com/avareum/avareum-hubble-signer/internal/types"
	"github.com/avareum/avareum-hubble-signer/pkg/logger"
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
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "OK",
			})
		})

		v1.POST("/:chain/:cluster/wallet/new", func(c *gin.Context) {
			res, err := walletHandler.NewWallet(wallet.NewWalletRequest{
				Chain: types.NewChain(c.Param("chain"), c.Param("cluster")),
			})
			if err != nil {
				logger.Default.Err(err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "FAILED",
					"error":  "create new wallet failed",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status":   "OK",
				"response": res,
			})
		})

		v1.POST("/:chain/:cluster/wallet/execute", func(c *gin.Context) {
			// TODO: parse request body
			res, err := walletHandler.Execute(wallet.ExecuteWalletRequest{
				Chain: types.NewChain(c.Param("chain"), c.Param("cluster")),
			})
			if err != nil {
				logger.Default.Err(err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "FAILED",
					"error":  "execute wallet failed",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status":   "OK",
				"response": res,
			})
		})
	}
	r.Run()
}
