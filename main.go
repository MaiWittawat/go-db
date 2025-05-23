package main

import (
	"context"
	"go-rebuild/db"
	"go-rebuild/handler/api"
	handler "go-rebuild/handler/user"
	"go-rebuild/module/port"
	"go-rebuild/module/service"
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

	var userDB port.UserDB
	useMongo := true

	if useMongo {
		mgDB, err := db.InitMongoDB(ctx)
		if err != nil {
			log.Fatal("fail to connect mongodb: ", err)
		}
		userDB = repository.NewMongoUserRepo(mgDB)
	} else {
		pgDB, err := db.InitPsqlDB()
		if err != nil {
			log.Fatal("fail to connect psqldb: ", err)
		}
		userDB = repository.NewPsqlUserRepo(pgDB)
	}

	router := gin.Default()
	
	userRepo := repository.NewUserRepo(userDB)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)
	api.RegisterUserAPI(router, userHandler)
	
	
	if err := router.Run(":3000"); err != nil{
		log.Fatal("fatil to start server")
	}	
	log.Println("start server at port 3000")
}