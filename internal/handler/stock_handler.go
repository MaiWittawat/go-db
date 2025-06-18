package handler

import (
	"go-rebuild/internal/module"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	service module.StockService
}

func NewStockHandler(service module.StockService) *StockHandler {
	return &StockHandler{
		service: service,
	}
}

func (h *StockHandler) GetStocks(c *gin.Context) {
	stocks, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "get stocks success", "data": stocks})
}

func (h *StockHandler) GetStock(c *gin.Context) {
	stock, err := h.service.GetByProductID(c.Request.Context(), c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "get stock success", "data": stock})
}
