package handler

import (
	"go-rebuild/internal/auth"
	"go-rebuild/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	service auth.Jwt
}

func NewAuthHandler(service auth.Jwt) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var user model.User
	user.Role = "USER"

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	if err := h.service.Register(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}


func (h *AuthHandler) RegisterSeller(c *gin.Context) {
	var seller model.User
	seller.Role = "SELLER"

	if err := c.ShouldBindJSON(&seller); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	if err := h.service.Register(c.Request.Context(), &seller); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	logrus.Info("Login")
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}