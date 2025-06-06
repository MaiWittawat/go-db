package handler

import (
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service module.ProductService
}

func NewProductHandler(service module.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Save(c.Request.Context(), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "product created"})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var upDateProduct model.Product
	if err := c.ShouldBindJSON(&upDateProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Update(c.Request.Context(), &upDateProduct, c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product updated"})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "get product success", "data": products})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	var product model.Product
	if err := h.service.GetByID(c.Request.Context(), c.Param("id"), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "get product success", "data": product})
}
