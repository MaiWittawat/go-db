package api

import (
	handler "go-rebuild/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserAPI(router *gin.Engine, userHandler *handler.UserHandler) {
	protected := router.Group("/users")
	protected.POST("/register/user", userHandler.RegisterUser)
	protected.POST("/register/seller", userHandler.RegisterSeller)
	protected.GET("/", userHandler.GetUsers)
	protected.GET("/:id", userHandler.GetUserByID)
	protected.PATCH("/:id", userHandler.EditUser)
	protected.DELETE("/:id", userHandler.DropUser)
}