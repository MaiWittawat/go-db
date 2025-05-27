package handler

import (
	"fmt"
	"go-rebuild/model"
	module "go-rebuild/module/user"
	"net/http"

	"github.com/gin-gonic/gin"
)


type UserHandler struct {
	service module.UserService
}

func NewUserHandler(service module.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (uh *UserHandler) RegisterUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uh.service.SaveUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (uh *UserHandler) RegisterSeller(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uh.service.SaveSeller(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "seller created"})
}

func (uh *UserHandler) EditUser(c *gin.Context) {
	var user model.User
	id := c.Param("id")
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uh.service.Update(c.Request.Context(), &user, id)
	if err != nil {
		fmt.Println("error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (uh *UserHandler) DropUser(c *gin.Context) {
	id := c.Param("id")

	err := uh.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}



func (uh *UserHandler) GetUsers(c *gin.Context) {
	users, err := uh.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) 
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "get all user success", "data": users})
}


func (uh *UserHandler) GetUserByID(c *gin.Context) {
	user, err := uh.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) 
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "get all user success", "data": user})
}
