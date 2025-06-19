package api

import (
	"go-rebuild/internal/auth"
	"go-rebuild/internal/handler"

	"github.com/gin-gonic/gin"
)

func AuthAPI(router *gin.Engine, authHandler *handler.AuthHandler) {
	router.POST("/register/user", authHandler.RegisterUser)
	router.POST("/register/seller", authHandler.RegisterSeller)
	router.POST("/login", authHandler.Login)
}

func UserAPI(router *gin.Engine, userHandler *handler.UserHandler, authSvc auth.Jwt) {
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

func OrderAPI(router *gin.Engine, orderHandler *handler.OrderHandler, authSvc auth.Jwt) {
	protected := router.Group("/orders")
	protected.Use(
		handler.AuthenticateMiddleware(authSvc),
		handler.AuthorizeMiddleware(authSvc, "USER", "SELLER", "ADMIN"),
	)
	protected.GET("/:id", orderHandler.GetOrder)
	protected.POST("/", orderHandler.CreateOrder)
	protected.PATCH("/:id", orderHandler.UpdateOrder)
	protected.DELETE("/:id", orderHandler.DeleteOrder)

	adminOnly := router.Group("/orders").Use(handler.AuthorizeMiddleware(authSvc, "ADMIN"))
	adminOnly.GET("/", orderHandler.GetOrders)
}

func ProductAPI(router *gin.Engine, productHandler *handler.ProductHandler, authSvc auth.Jwt) {
	public := router.Group("/products")
	public.GET("/", productHandler.GetProducts)
	public.GET("/:id", productHandler.GetProduct)

	protected := router.Group("/products")
	protected.Use(
		handler.AuthenticateMiddleware(authSvc),
		handler.AuthorizeMiddleware(authSvc, "SELLER", "ADMIN"),
	)
	protected.POST("/", productHandler.CreateProduct)
	protected.PATCH("/:id", productHandler.UpdateProduct)
	protected.DELETE("/:id", productHandler.DeleteProduct)
}

func StockAPI(router *gin.Engine, stockHandler *handler.StockHandler, authSvc auth.Jwt) {
	public := router.Group("/stocks")
	public.GET("/", stockHandler.GetStocks)
	public.GET("/:product_id", stockHandler.GetStock)
}

func MessageAPI(router *gin.Engine, messageHandler *handler.MessageHandler, authSvc auth.Jwt) {
	protected := router.Group("/messages")
	protected.GET("/ws", messageHandler.Connect) // realtime chat
	protected.GET("/:user_id1/user_id2", messageHandler.GetMessagesBetweenUser)
}

func StorageAPI(router *gin.Engine, storerHandler *handler.StorerHandler) {
	public := router.Group("/")
	public.GET("/download", storerHandler.GetFileUrl)
	public.POST("/upload", storerHandler.Upload)
}
