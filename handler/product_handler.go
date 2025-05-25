package handler

import (
	"go-rebuild/model"
	module "go-rebuild/module/product"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service module.ProductService
}

func NewProductHandler(service module.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (ph *ProductHandler) CreateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ph.service.Save(c.Request.Context(), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "product created"})
}

func (ph *ProductHandler) UpdateProduct(c *gin.Context) {
	var upDateProduct model.Product
	if err := c.ShouldBindJSON(&upDateProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ph.service.Update(c.Request.Context(), &upDateProduct, c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return	
	}

	c.JSON(http.StatusCreated, gin.H{"message": "product updated"})
}

func (ph *ProductHandler) DeleteProduct(c *gin.Context) {
	if err := ph.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "product deleted"})
}