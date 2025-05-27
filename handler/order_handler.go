package handler

import (
	"go-rebuild/model"
	"net/http"
	module "go-rebuild/module/order"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service module.OrderService
}

func NewOrderHandler(service module.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (ph *OrderHandler) CreateOrder(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ph.service.Save(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "order created"})
}

func (ph *OrderHandler) UpdateOrder(c *gin.Context) {
	var upDateOrder model.Order
	if err := c.ShouldBindJSON(&upDateOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ph.service.Update(c.Request.Context(), &upDateOrder, c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return	
	}

	c.JSON(http.StatusOK, gin.H{"message": "order updated"})
}

func (ph *OrderHandler) DeleteOrder(c *gin.Context) {
	if err := ph.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}


func (ph *OrderHandler) GetOrder(c *gin.Context) {
	var order model.Order
	if err := ph.service.GetByID(c.Request.Context(), c.Param("id"), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "get order success", "data": order})
}


func (ph *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := ph.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "get orders success", "data": orders})
}
