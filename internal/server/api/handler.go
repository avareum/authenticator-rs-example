package api

import (
	"encoding/base64"
	"net/http"

	"github.com/avareum/avareum-hubble-signer/internal/server/wallet"
	"github.com/avareum/avareum-hubble-signer/internal/types"
	"github.com/avareum/avareum-hubble-signer/pkg/logger"
	"github.com/gin-gonic/gin"

	signertypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
)

type IRestHandler interface {
	NewWallet(c *gin.Context)
	Execute(c *gin.Context)
}
type RestHandler struct {
	walletHandler wallet.WalletHandler
}

func NewRestHandler() IRestHandler {
	r := &RestHandler{
		walletHandler: wallet.NewFundWalletHandler(),
	}
	return r
}

// @BasePath 	/v1
// @Summary 	Create a new wallet
// @Description Create & store a new wallet of the given chain and cluster
// @Accept 		mpfd
// @Produce 	json
// @Param 		chain path string required "Chain name"
// @Param 		chain_id path string required "Chain id"
// @Success 	200 {object} wallet.NewWalletResponse
// @Router 		/{chain}/{chain_id}/wallet/new [post]
func (r *RestHandler) NewWallet(c *gin.Context) {
	res, err := r.walletHandler.NewWallet(wallet.NewWalletRequest{
		Chain: types.NewChain(c.Param("chain"), c.Param("chain_id")),
	})
	if err != nil {
		logger.Default.Err(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "FAILED",
			"error":  "create new wallet failed",
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

// @BasePath 	/v1
// @Summary 	Execute a wallet
// @Description Execute a payload with the given wallet
// @Accept  	mpfd
// @Produce  	json
// @Param 		chain path string required "Target chain name"
// @Param 		chain_id path string required "Target chain id"
// @Param 		wallet path string required "Target fund wallet address"
// @Param 		caller formData string required "Caller service name"
// @Param 		payload formData string required "Transaction payload"
// @Param 		signature formData string required "Tranasction payload signature"
// @Success 	200 {object} wallet.ExecuteWalletResponse
// @Router 	/{chain}/{chain_id}/{wallet}/execute [post]
func (r *RestHandler) Execute(c *gin.Context) {
	chain := types.NewChain(c.Param("chain"), c.Param("chain_id"))
	payload, err := base64.StdEncoding.DecodeString(c.PostForm("payload"))
	if err != nil {
		logger.Default.Err(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "FAILED",
			"error":  "decode payload failed",
		})
		return
	}
	signature, err := base64.StdEncoding.DecodeString(c.PostForm("signature"))
	if err != nil {
		logger.Default.Err(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "FAILED",
			"error":  "decode payload signature failed",
		})
		return
	}
	request := wallet.ExecuteWalletRequest{
		Chain: chain,
		SignerRequest: signertypes.SignerRequest{
			Chain:     chain,
			Wallet:    c.Param("wallet"),
			Caller:    c.PostForm("caller"),
			Payload:   payload,
			Signature: signature,
		},
	}
	res, err := r.walletHandler.Execute(request)
	if err != nil {
		logger.Default.Err(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "FAILED",
			"error":  "execute wallet failed",
		})
		return
	}
	c.JSON(http.StatusOK, res)
}
