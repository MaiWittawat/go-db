package api

import (
	"go-rebuild/internal/auth"
	"go-rebuild/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterOrderAPI(router *gin.Engine, orderHandler *handler.OrderHandler, authSvc auth.Jwt) {

	protected := router.Group("/orders")
	protected.Use(
		handler.AuthenticateMiddleware(authSvc),
		handler.AuthorizeMiddleware(authSvc, "USER", "SELLER", "ADMIN"),
	)
	protected.GET("/:id", orderHandler.GetOrder)
	protected.POST("/", orderHandler.CreateOrder)
	protected.PATCH("/:id", orderHandler.UpdateOrder)
	protected.DELETE("/:id", orderHandler.DeleteOrder)

	adminOnly := router.Group("/orders").Use(handler.AuthorizeMiddleware(authSvc, "ADMIN"))
	adminOnly.GET("/", orderHandler.GetOrders)
}
