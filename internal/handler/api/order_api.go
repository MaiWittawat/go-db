package api

import (
	"go-rebuild/internal/auth"
	"go-rebuild/internal/handler"
	"go-rebuild/internal/handler/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterOrderAPI(router *gin.Engine, orderHandler *handler.OrderHandler, authSvc auth.AuthService) {

	protected := router.Group("/orders").Use(middleware.AuthorizeMiddleware(authSvc, "USER", "SELLER", "ADMIN"),)
	protected.GET("/:id", orderHandler.GetOrder)
	protected.POST("/", orderHandler.CreateOrder)
	protected.PATCH("/:id", orderHandler.UpdateOrder)
	protected.DELETE("/:id", orderHandler.DeleteOrder)

	adminOnly := router.Group("/orders").Use(middleware.AuthorizeMiddleware(authSvc, "ADMIN",))
	adminOnly.GET("/", orderHandler.GetOrders)
}