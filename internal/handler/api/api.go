package api

import (
	"go-rebuild/internal/auth"
	"go-rebuild/internal/handler"

	"github.com/gin-gonic/gin"
)

type APIRouterConfigurator struct {
	authHandler    *handler.AuthHandler
	userHandler    *handler.UserHandler
	orderHandler   *handler.OrderHandler
	productHandler *handler.ProductHandler
	stockHandler   *handler.StockHandler
	messageHandler *handler.MessageHandler
	storerHandler  *handler.StorerHandler
	auth           auth.Jwt
}

func NewAPIRouterConfigurator(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, orderHandler *handler.OrderHandler, productHandler *handler.ProductHandler, stockHandler *handler.StockHandler, messageHandler *handler.MessageHandler, storerHandler *handler.StorerHandler, auth auth.Jwt) *APIRouterConfigurator {
	return &APIRouterConfigurator{
		authHandler:    authHandler,
		userHandler:    userHandler,
		orderHandler:   orderHandler,
		productHandler: productHandler,
		stockHandler:   stockHandler,
		messageHandler: messageHandler,
		storerHandler:  storerHandler,
		auth:           auth,
	}
}

func (api *APIRouterConfigurator) PublicAPIRoutes(router *gin.Engine) {
	public := router.Group("/")
	public.POST("/register/user", api.authHandler.RegisterUser)
	public.POST("/register/seller", api.authHandler.RegisterSeller)
	public.POST("/login", api.authHandler.Login)

	// Product
	publicProduct := router.Group("/products")
	publicProduct.GET("/", api.productHandler.GetProducts)
	publicProduct.GET("/:id", api.productHandler.GetProduct)

	// Stock
	publicStock := router.Group("/stocks")
	publicStock.GET("/", api.stockHandler.GetStocks)
	publicStock.GET("/:product_id", api.stockHandler.GetStock)

	// Storer
	publicStorer := router.Group("/")
	publicStorer.GET("/download", api.storerHandler.GetFileUrl)
	publicStorer.POST("/upload", api.storerHandler.Upload)

	// Message
	publicMessage := router.Group("/messages")
	publicMessage.GET("/ws", api.messageHandler.Connect) // realtime chat
	publicMessage.GET("/:user_id1/user_id2", api.messageHandler.GetMessagesBetweenUser)
}

func (api *APIRouterConfigurator) ProtectAPIRoutes(router *gin.Engine) {
	// User
	protectedUser := router.Group("/users")
	protectedUser.Use(handler.AuthenticateMiddleware(api.auth), handler.AuthorizeMiddleware(api.auth, "USER", "SELLER", "ADMIN"))
	protectedUser.GET("/", api.userHandler.GetUsers)
	protectedUser.GET("/:id", api.userHandler.GetUserByID)
	protectedUser.PATCH("/:id", api.userHandler.EditUser)
	protectedUser.DELETE("/:id", api.userHandler.DropUser)

	// Product
	protectedProduct := router.Group("/products")
	protectedProduct.Use(handler.AuthenticateMiddleware(api.auth), handler.AuthorizeMiddleware(api.auth, "SELLER", "ADMIN"))
	protectedProduct.POST("/", api.productHandler.CreateProduct)
	protectedProduct.PATCH("/:id", api.productHandler.UpdateProduct)
	protectedProduct.DELETE("/:id", api.productHandler.DeleteProduct)

	// Order
	protectedOrder := router.Group("/orders")
	protectedOrder.Use(handler.AuthenticateMiddleware(api.auth), handler.AuthorizeMiddleware(api.auth, "USER", "SELLER", "ADMIN"))
	protectedOrder.GET("/:id", api.orderHandler.GetOrder)
	protectedOrder.POST("/", api.orderHandler.CreateOrder)
	protectedOrder.PATCH("/:id", api.orderHandler.UpdateOrder)
	protectedOrder.DELETE("/:id", api.orderHandler.DeleteOrder)
	protectedOrder.GET("/", api.orderHandler.GetOrders).Use(handler.AuthorizeMiddleware(api.auth, "ADMIN"))
}
