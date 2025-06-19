package handler

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	RegisterUser(c *gin.Context)
	RegisterSeller(c *gin.Context)
	Login(c *gin.Context)
}

type UserHandler interface {
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)

	GetUsers(c *gin.Context)
	GetUserByID(c *gin.Context)
}

type ProductHandler interface {
	CreateProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)

	GetProducts(c *gin.Context)
	GetProductByID(c *gin.Context)
}

type OrderHandler interface {
	CreateOrder(c *gin.Context)
	UpdateOrder(c *gin.Context)
	DeleteOrder(c *gin.Context)

	GetOrders(c *gin.Context)
	GetOrderByID(c *gin.Context)
}

type StockHandler interface {
	GetStocks(c *gin.Context)
	GetStockByProductID(c *gin.Context)
}

type MessageHandler interface {
	Connect(c *gin.Context)
	GetMessagesBetweenUser(c *gin.Context)
}

type StorerHandler interface {
	GetFileUrl(c *gin.Context)
	Upload(c *gin.Context)
}