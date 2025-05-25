package api

import (
	handler "go-rebuild/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserAPI(router *gin.Engine, userHandler *handler.UserHandler) {
	public := router.Group("/users")

	public.POST("/", userHandler.RegisterUser)
	public.PATCH("/:id", userHandler.EditUser)
	public.DELETE("/:id", userHandler.DropUser)
}