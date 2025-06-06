package api

import (
	"go-rebuild/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAuthAPI(router *gin.Engine, authHandler handler.AuthHandler) {
	router.POST("/register/user", authHandler.RegisterUser)
	router.POST("/register/seller", authHandler.RegisterSeller)
	router.POST("/login", authHandler.Login)
}