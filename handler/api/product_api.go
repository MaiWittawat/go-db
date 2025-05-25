package api

import (
	handler "go-rebuild/handler"

	"github.com/gin-gonic/gin"
)

func RegisterProductAPI(router *gin.Engine, productHandler *handler.ProductHandler) {
	public := router.Group("/products")

	public.POST("/", productHandler.CreateProduct)
	public.PATCH("/:id", productHandler.UpdateProduct)
	public.DELETE("/:id", productHandler.DeleteProduct)
}