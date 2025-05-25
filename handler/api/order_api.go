package api

import (
	"go-rebuild/handler"

	"github.com/gin-gonic/gin"
)

func RegisterOrderAPI(router *gin.Engine, orderHandler *handler.OrderHandler) {
	public := router.Group("/orders")

	public.POST("/", orderHandler.CreateOrder)
	public.PATCH("/:id", orderHandler.UpdateOrder)
	public.DELETE("/:id", orderHandler.DeleteOrder)
}