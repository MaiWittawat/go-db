package main

// import (
// 	"go-rebuild/auth"
// 	dbRepo "go-rebuild/db"
// 	"go-rebuild/handler"
// 	moduleOrder "go-rebuild/module/order"
// 	moduleProduct "go-rebuild/module/product"
// 	moduleUser "go-rebuild/module/user"
// 	rclient "go-rebuild/redis"
// 	"go-rebuild/repository"

// 	"github.com/gin-gonic/gin"
// 	"github.com/redis/go-redis/v9"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"gorm.io/gorm"
// )

// type Container struct {
// 	// Database connections
// 	MgDBClient  *mongo.Client
// 	PgDB        *gorm.DB
// 	RedisClient *redis.Client

// 	// Core services
// 	DBRepo   dbRepo.DB
// 	CacheSvc rclient.Cache
// 	AuthRepo auth.AuthRepo
// 	Router   *gin.Engine

// 	// Repositories
// 	UserRepo    repository.UserRepo
// 	ProductRepo repository.ProductRepo
// 	OrderRepo   repository.OrderRepo

// 	// Services
// 	UserSvc    moduleUser.UserService
// 	ProductSvc moduleProduct.ProductService
// 	OrderSvc   moduleOrder.OrderService
// 	AuthSvc    auth.AuthService

// 	// Handlers
// 	UserHandler    handler.UserHandler
// 	AuthHandler    handler.AuthHandler
// 	ProductHandler handler.ProductHandler
// 	OrderHandler   handler.OrderHandler
// }
