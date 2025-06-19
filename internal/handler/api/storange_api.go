package api

import (
	"go-rebuild/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterStorageAPI(router *gin.Engine, storerHandler *handler.StorerHandler) {
	public := router.Group("/")
	public.GET("/download", storerHandler.GetFileUrl)
	public.POST("/upload", storerHandler.Upload)
}