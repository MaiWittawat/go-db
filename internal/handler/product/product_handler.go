package product

import (
	"go-rebuild/internal/handler"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"net/http"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	service module.ProductService
}

func NewProductHandler(service module.ProductService) handler.ProductHandler {
	return &productHandler{service: service}
}

func (h *productHandler) CreateProduct(c *gin.Context) {
	var productReq model.ProductReq
	
	if err := c.ShouldBindJSON(&productReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	userID := c.GetString("user_id")

	if err := h.service.Save(c.Request.Context(), &productReq, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "product created"})
}

func (h *productHandler) UpdateProduct(c *gin.Context) {
	var upDateProductReq model.ProductReq
	userID := c.GetString("user_id")

	if err := c.ShouldBindJSON(&upDateProductReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), &upDateProductReq, c.Param("id"), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product updated"})
}

func (h *productHandler) DeleteProduct(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}

func (h *productHandler) GetProducts(c *gin.Context) {
	productsRes, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "get product success", "data": productsRes})
}

func (h *productHandler) GetProductByID(c *gin.Context) {
	productRes, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "get product success", "data": productRes})
}
