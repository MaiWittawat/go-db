package handler

import (
	"go-rebuild/model"
	"go-rebuild/module/port"

	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	Username string `json:"username"`
	Password string	`json:"password"`
	Email string	`json:"email"`
}

type UserHandler struct {
	service port.UserService
}

func NewUserHandler(service port.UserService) *UserHandler {
	return &UserHandler{service: service}
}


func request2User(userReq *UserRequest) *model.User {
	return &model.User{
		Username: userReq.Username,
		Password: userReq.Password,
		Email: userReq.Email,
	}
}

func (uh *UserHandler)	RegisterUser(c *gin.Context) {
	var userReq UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := uh.service.Save(c.Request.Context(), request2User(&userReq))
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot save user"})
		return
	}

	c.JSON(201, gin.H{"message": "user created"})
}

func (uh *UserHandler)	EditUser(c *gin.Context) {
	var userReq UserRequest
	id := c.Param("id")
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := uh.service.Update(c.Request.Context(), request2User(&userReq), id)
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot update user"})
		return
	}

	c.JSON(200, gin.H{"message": "user updated"})
}

func (uh *UserHandler) DropUser(c *gin.Context) {
	id := c.Param("id")

	err := uh.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot delete user"})
		return
	}

	c.JSON(200, gin.H{"message": "user deleted"})
}