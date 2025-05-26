package main

import (
	"context"
	"go-rebuild/db"
	"go-rebuild/handler/api"
	"go-rebuild/handler"
	moduleUser "go-rebuild/module/user"
	moduleProduct "go-rebuild/module/product"
	moduleOrder "go-rebuild/module/order"
	"go-rebuild/repository"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main(){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	godotenv.Load()

	var dbRepo db.DB
	useMongo := false

	if useMongo {
		mgDB, err := db.InitMongoDB(ctx)
		if err != nil {
			log.Fatal("fail to connect mongodb: ", err)
		}
		log.Println("connect to mongo success")
		dbRepo = db.NewMongoRepo(mgDB, "miniproject")
	} else {
		pgDB, err := db.InitPsqlDB()
		if err != nil {
			log.Fatal("fail to connect psqldb: ", err)
		}
		log.Println("connect to psql success")
		dbRepo = db.NewPsqlRepo(pgDB)
	}

	router := gin.Default()
	
	userRepo := repository.NewUserRepo(dbRepo)
	userSvc := moduleUser.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)
	api.RegisterUserAPI(router, userHandler)


	productRepo := repository.NewProductRepo(dbRepo)
	productSvc := moduleProduct.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)
	api.RegisterProductAPI(router, productHandler)

	orderRepo := repository.NewOrderRepo(dbRepo)
	orderSvc := moduleOrder.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderSvc)
	api.RegisterOrderAPI(router, orderHandler)

	
	
	if err := router.Run(":3000"); err != nil{
		log.Fatal("fatil to start server")
	}	
	log.Println("start server at port 3000")
}