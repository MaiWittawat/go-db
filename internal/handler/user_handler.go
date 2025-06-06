package handler

import (
	"fmt"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service module.UserService
}

func NewUserHandler(service module.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) EditUser(c *gin.Context) {
	verifyId := c.GetString("user_id")
	if verifyId != c.Param("id") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden id not match"})
		return
	}
	var user model.User
	id := c.Param("id")
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Update(c.Request.Context(), &user, id)
	if err != nil {
		fmt.Println("error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (h *UserHandler) DropUser(c *gin.Context) {
	verifyId := c.GetString("user_id")
	if verifyId != c.Param("id") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden id not match"})
		return
	}

	err := h.service.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "get all user success", "data": users})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	user, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "get all user success", "data": user})
}
