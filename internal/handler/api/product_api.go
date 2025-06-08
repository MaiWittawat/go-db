package api

import (
	"go-rebuild/internal/auth"
	handler "go-rebuild/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterProductAPI(router *gin.Engine, productHandler *handler.ProductHandler, authSvc auth.Jwt) {
	public := router.Group("/products")
	public.GET("/", productHandler.GetProducts)
	public.GET("/:id", productHandler.GetProduct)

	protected := router.Group("/products")
	protected.Use(
		handler.AuthenticateMiddleware(authSvc),
		handler.AuthorizeMiddleware(authSvc, "SELLER", "ADMIN"),
	)
	protected.POST("/", productHandler.CreateProduct)
	protected.PATCH("/:id", productHandler.UpdateProduct)
	protected.DELETE("/:id", productHandler.DeleteProduct)
}