package api

import (
	"go-rebuild/internal/auth"
	handler "go-rebuild/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterUserAPI(router *gin.Engine, userHandler *handler.UserHandler, authSvc auth.Jwt) {
	protected := router.Group("/users")
	protected.Use(
		handler.AuthenticateMiddleware(authSvc),
		handler.AuthorizeMiddleware(authSvc, "USER", "SELLER", "ADMIN"),
	)

	protected.GET("/", userHandler.GetUsers)
	protected.GET("/:id", userHandler.GetUserByID)
	protected.PATCH("/:id", userHandler.EditUser)
	protected.DELETE("/:id", userHandler.DropUser)
}