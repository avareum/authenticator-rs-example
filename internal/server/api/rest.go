package api

import (
	"net/http"

	"github.com/avareum/avareum-hubble-signer/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RestAPI struct{}

func NewRestAPI() *RestAPI {
	return &RestAPI{}
}

func (api *RestAPI) Serve() {
	handler := NewRestHandler()
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "OK",
			})
		})
		v1.POST("/:chain/:chain_id/wallet/new", handler.NewWallet)
		v1.POST("/:chain/:chain_id/:wallet/execute", handler.Execute)
	}

	docs.SwaggerInfo.BasePath = "/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Run(":8080")
}
