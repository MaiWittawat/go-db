package api

import (
	"go-rebuild/internal/auth"
	"go-rebuild/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterMessageAPI(router *gin.Engine, messageHandler *handler.MessageHandler, authSvc auth.Jwt) {
	protected := router.Group("/messages")
	protected.GET("/ws", messageHandler.Connect) // realtime chat

	// protected.Use(
	// 	handler.AuthenticateMiddleware(authSvc),
	// 	handler.AuthorizeMiddleware(authSvc, "USER", "SELLER", "ADMIN"),
	// )
	protected.GET("/:user_id1/user_id2", messageHandler.GetMessagesBetweenUser)
}
