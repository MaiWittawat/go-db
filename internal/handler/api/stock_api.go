package api

import (
	"go-rebuild/internal/auth"
	"go-rebuild/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterStockAPI(router *gin.Engine, stockHandler *handler.StockHandler, authSvc auth.Jwt) {
	public := router.Group("/stocks")
	public.GET("/", stockHandler.GetStocks)
	public.GET("/:product_id", stockHandler.GetStock)
}