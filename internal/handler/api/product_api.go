package api

import (
	"go-rebuild/internal/auth"
	handler "go-rebuild/internal/handler"
	"go-rebuild/internal/handler/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterProductAPI(router *gin.Engine, productHandler *handler.ProductHandler, authSvc auth.AuthService) {
	public := router.Group("/products")
	public.GET("/", productHandler.GetProducts)
	public.GET("/:id", productHandler.GetProduct)

	protected := router.Group("/products")
	protected.Use(
		middleware.JWTAuthenMiddleware(authSvc),
		middleware.AuthorizeMiddleware(authSvc, "SELLER", "ADMIN"),
	)
	protected.POST("/", productHandler.CreateProduct)
	protected.PATCH("/:id", productHandler.UpdateProduct)
	protected.DELETE("/:id", productHandler.DeleteProduct)
}