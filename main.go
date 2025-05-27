package main

import (
	"context"
	"go-rebuild/db"
	"go-rebuild/handler"
	"go-rebuild/handler/api"
	moduleOrder "go-rebuild/module/order"
	moduleProduct "go-rebuild/module/product"
	moduleUser "go-rebuild/module/user"
	"go-rebuild/redis"
	"go-rebuild/repository"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	godotenv.Load()

	if os.Getenv("ENV") == "develop" {
		log.SetLevel(log.TraceLevel)
	} 

	var dbRepo db.DB
	useMongo := false

	if useMongo {
		mgDB, err := db.InitMongoDB(ctx)
		if err != nil {
			log.Panic("fail to connect mongodb: ", err)
		}
		log.Info("connect to mongo success")
		dbRepo = db.NewMongoRepo(mgDB, "miniproject")
	} else {
		pgDB, err := db.InitPsqlDB()
		if err != nil {
			log.Panic("fail to connect psqldb: ", err)
		}
		log.Info("connect to psql success")
		dbRepo = db.NewPsqlRepo(pgDB)
	}

	redisClient := redis.InitRedisClient()
	redisCache := redis.NewRedisCache(redisClient)
	router := gin.Default()

	userRepo := repository.NewUserRepo(dbRepo, redisCache)
	userSvc := moduleUser.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)
	api.RegisterUserAPI(router, userHandler)

	productRepo := repository.NewProductRepo(dbRepo, redisCache)
	productSvc := moduleProduct.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)
	api.RegisterProductAPI(router, productHandler)

	orderRepo := repository.NewOrderRepo(dbRepo, redisCache)
	orderSvc := moduleOrder.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderSvc)
	api.RegisterOrderAPI(router, orderHandler)

	if err := router.Run(":3000"); err != nil {
		log.Panic("fatil to start server: ", err)
	}

	log.Info("server close ...")
}
