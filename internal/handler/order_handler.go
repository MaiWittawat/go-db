package handler

import (
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service module.OrderService
}

func NewOrderHandler(service module.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetString("user_id")
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.UserID = userID
	if err := h.service.Save(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "order created"})
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	var upDateOrder model.Order
	if err := c.ShouldBindJSON(&upDateOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Update(c.Request.Context(), &upDateOrder, c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return	
	}

	c.JSON(http.StatusOK, gin.H{"message": "order updated"})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	userID := c.GetString("user_id")
	if err := h.service.Delete(c.Request.Context(), c.Param("id"), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}


func (h *OrderHandler) GetOrder(c *gin.Context) {
	var order model.Order
	if err := h.service.GetByID(c.Request.Context(), c.Param("id"), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "get order success", "data": order})
}


func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "get orders success", "data": orders})
}
