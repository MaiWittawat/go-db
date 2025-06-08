package main

import (
	"context"
	appcore_config "go-rebuild/cmd/go-rebuild/config"
	"go-rebuild/internal/auth"
	rclient "go-rebuild/internal/cache"
	"go-rebuild/internal/db"
	"go-rebuild/internal/handler"
	"go-rebuild/internal/handler/api"

	orderSvc "go-rebuild/internal/module/order"
	productSvc "go-rebuild/internal/module/product"
	userSvc "go-rebuild/internal/module/user"
	orderRepo "go-rebuild/internal/repository/order"
	productRepo "go-rebuild/internal/repository/product"
	userRepo "go-rebuild/internal/repository/user"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	mgDBInstant        *mongo.Client
	pgDBInstant        *gorm.DB
	redisClientInstant *redis.Client
)

func main() {
	appcore_config.InitConfigurations()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if appcore_config.Config.Mode == "develop" {
		log.SetLevel(log.InfoLevel)
	}

	var dbRepo db.DB
	useMongo := false

	if useMongo {
		initMongoCtx, initMongoCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer initMongoCancel()

		mgClient, err := db.InitMongoDB(initMongoCtx)
		if err != nil {
			log.Panic("fail to connect mongodb: ", err)
		}

		log.Info("connect to mongo success")
		mgDBInstant = mgClient
		dbRepo = db.NewMongoRepo(mgDBInstant, "miniproject")

	} else {
		pg, err := db.InitPsqlDB()

		if err != nil {
			log.Panic("fail to connect psqldb: ", err)
		}

		log.Info("connect to psql success")
		pgDBInstant = pg
		dbRepo = db.NewPsqlRepo(pgDBInstant)
	}

	redisClient := rclient.InitRedisClient()
	cacheSvc := rclient.NewCacheService(redisClient)
	router := gin.Default()

	
	// User and Auth
	userRepository := userRepo.NewUserRepo(dbRepo, cacheSvc)
	authService := auth.NewAuth(userRepository)

	userService := userSvc.NewUserService(userRepository)

	userHandler := handler.NewUserHandler(userService)
	api.RegisterUserAPI(router, userHandler, authService)

	authHandler := handler.NewAuthHandler(authService)
	api.RegisterAuthAPI(router, authHandler)

	// Product
	productRepo := productRepo.NewProductRepo(dbRepo, cacheSvc)
	productSvc := productSvc.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)
	api.RegisterProductAPI(router, productHandler, authService)

	// Order
	orderRepo := orderRepo.NewOrderRepo(dbRepo, cacheSvc)
	orderSvc := orderSvc.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderSvc)
	api.RegisterOrderAPI(router, orderHandler, authService)

	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Info("Server closed gracefully")
			} else {
				log.Fatalf("listen: %s\n", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown signal received")

	gracefulShutdown(shutdownCtx, server)

}

func gracefulShutdown(ctx context.Context, server *http.Server) {
	log.Info("Shutting down server...")

	// close HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("HTTP server Shutdown: %v", err)
	}

	// close Redis
	if redisClientInstant != nil {
		if err := redisClientInstant.Close(); err != nil {
			log.Errorf("Redis shutdown error: %v", err)
		} else {
			log.Info("Redis closed")
		}
	}

	// clse MongoDB
	if mgDBInstant != nil {
		if err := mgDBInstant.Disconnect(ctx); err != nil {
			log.Errorf("MongoDB shutdown error: %v", err)
		} else {
			log.Info("MongoDB disconnected")
		}
	}

	// close Postgres
	if pgDBInstant != nil {
		sqlDB, err := pgDBInstant.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Errorf("Postgres shutdown error: %v", err)
			} else {
				log.Info("Postgres closed")
			}
		}
	}

	log.Info("Graceful shutdown complete.")
}
