package main

import (
	"context"
	appcore_config "go-rebuild/cmd/go-rebuild/config"
	"go-rebuild/internal/auth"
	redisclient "go-rebuild/internal/cache"
	"go-rebuild/internal/db"
	"go-rebuild/internal/handler"
	"go-rebuild/internal/handler/api"
	"go-rebuild/internal/mail"
	messagebroker "go-rebuild/internal/message_broker"
	"go-rebuild/internal/model"
	"go-rebuild/internal/realtime"
	"go-rebuild/internal/storer"

	messageSvc "go-rebuild/internal/module/message"
	orderSvc "go-rebuild/internal/module/order"
	productSvc "go-rebuild/internal/module/product"
	stockSvc "go-rebuild/internal/module/stock"
	userSvc "go-rebuild/internal/module/user"
	messageRepo "go-rebuild/internal/repository/message"
	orderRepo "go-rebuild/internal/repository/order"
	productRepo "go-rebuild/internal/repository/product"
	stockRepo "go-rebuild/internal/repository/stock"
	userRepo "go-rebuild/internal/repository/user"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	mgDBInstant        *mongo.Client
	pgDBInstant        *gorm.DB
	redisClientInstant *redis.Client
	rabbitMQConn       *amqp.Connection
)

func main() {
	// ------------------------------ Setup Config ------------------------------
	appcore_config.InitConfigurations()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	gin.SetMode(gin.ReleaseMode)

	if appcore_config.Config.Mode == "develop" {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	// ------------------------------ Init db ------------------------------
	var dbRepo db.DB
	var err error
	useMongo := false

	if useMongo {
		initMongoCtx, initMongoCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer initMongoCancel()

		mgDBInstant, err = db.InitMongoDB(initMongoCtx)
		if err != nil {
			log.Panic("fail to connect mongodb: ", err)
		}

		dbRepo = db.NewMongoRepo(mgDBInstant, "miniproject")

	} else {
		pgDBInstant, err = db.InitPsqlDB()

		if err != nil {
			log.Panic("fail to connect psqldb: ", err)
		}

		dbRepo, err = db.NewPsqlRepo(pgDBInstant)
		if err != nil {
			log.Fatal(err)
		}
	}

	router := gin.Default()
	// ------------------------------ Init 3rd party ------------------------------
	// init redis client
	redisClientInstant = redisclient.InitRedisClient(appcore_config.Config.RedisUrl, appcore_config.Config.RedisPass)
	cacheSvc := redisclient.NewCacheService(redisClientInstant)

	// init mail client
	mailClient := mail.InitSMTP()
	mailService := mail.NewMailService(mailClient)

	// init websocket
	websocketServer := realtime.NewWebSocketServer()

	// init minio
	minioClient := storer.InitMinio(appcore_config.Config.MinioURL, appcore_config.Config.MinioAccessKey, appcore_config.Config.MinioSecretKey)

	// init rabbitmq
	rabbitMQConn = messagebroker.InitRabbitmq()
	producerChannel := messagebroker.OpenChannel(rabbitMQConn)
	userConsumeChannel := messagebroker.OpenChannel(rabbitMQConn)
	stockConsumeChannel := messagebroker.OpenChannel(rabbitMQConn)

	if err := messagebroker.SetupExchangeAndQueue(userConsumeChannel, &model.MQConfig{
		ExchangeName: messagebroker.UserExchangeName,
		ExchangeType: messagebroker.UserExchangeType,
		QueueName:    messagebroker.UserQueueName,
		RoutingKey:   "user.#",
	}); err != nil {
		log.Fatalf("Failed to setup user exchange and queue: %v", err)
	}

	if err := messagebroker.SetupExchangeAndQueue(stockConsumeChannel, &model.MQConfig{
		ExchangeName: messagebroker.StockExchangeName,
		ExchangeType: messagebroker.StockExchangeType,
		QueueName:    messagebroker.StockQueueName,
		RoutingKey:   "stock.#",
	}); err != nil {
		log.Fatalf("Failed to setup stock exchange and queue: %v", err)
	}

	// ------------------------------ Start service ------------------------------
	// Repository
	userRepository := userRepo.NewUserRepo(dbRepo, cacheSvc)
	ProductRepository := productRepo.NewProductRepo(dbRepo, cacheSvc)
	orderRepository := orderRepo.NewOrderRepo(dbRepo, cacheSvc)
	stockRepository := stockRepo.NewStockRepo(dbRepo, cacheSvc)
	messageRepository := messageRepo.NewMessageRepo(dbRepo, cacheSvc)

	// Service
	stockService := stockSvc.NewStockService(stockRepository)
	producerService := messagebroker.NewProducer(producerChannel)
	consumerService := messagebroker.NewConsumer(userConsumeChannel, stockConsumeChannel, mailService, stockService)
	mqBroker := messagebroker.NewMessageBroker(producerService, consumerService)
	userService := userSvc.NewUserService(userRepository, producerService)
	authService := auth.NewAuthService(userService, producerService)
	productSvc := productSvc.NewProductService(ProductRepository, producerService)
	orderService := orderSvc.NewOrderService(orderRepository, productSvc, producerService)
	messageService := messageSvc.NewMessageService(messageRepository)
	liveChat := realtime.NewLiveChat(websocketServer, messageService, authService)
	storageService := storer.NewStorerService(minioClient, appcore_config.Config.MinioBucketName)

	// Handler
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productSvc)
	orderHandler := handler.NewOrderHandler(orderService)
	stockHandler := handler.NewStockHandler(stockService)
	messageHandler := handler.NewMessageHandler(liveChat, messageService)
	storerHandler := handler.NewStorerHandler(storageService)

	// API
	api.RegisterAuthAPI(router, authHandler)
	api.RegisterUserAPI(router, userHandler, authService)
	api.RegisterProductAPI(router, productHandler, authService)
	api.RegisterOrderAPI(router, orderHandler, authService)
	api.RegisterStockAPI(router, stockHandler, authService)
	api.RegisterMessageAPI(router, messageHandler, authService)
	api.RegisterStorageAPI(router, storerHandler)

	// start consume
	go mqBroker.EmailConsuming(messagebroker.UserQueueName, "user_consume")
	go mqBroker.StockConsuming(messagebroker.StockQueueName, "stock_consume")

	// ------------------------------ Start server ------------------------------
	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Info("[server]: Server closed gracefully")
			} else {
				log.Fatalf("listen: %s\n", err)
			}
		}
	}()

	log.Info("[server]: server start at port:3000")

	// ------------------------------ Shutdown ------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("[Signal]: shutdown signal received")

	// call shutdown all service
	gracefulShutdown(shutdownCtx, server)

}

// ------------------------------ Shutdown function ------------------------------
func rabbitShutdown() {
	// // close rabbitmq connection
	if rabbitMQConn != nil {
		if err := rabbitMQConn.Close(); err != nil {
			log.Errorf("[RabbitMQ] connection shutdown error: %v", err)
		}
	}
	log.Info("[server]: Rabbitmq closed")
}

func redisShutdown() {
	// close Redis
	if redisClientInstant != nil {
		if err := redisClientInstant.Close(); err != nil {
			log.Errorf("[Redis] shutdown error: %v", err)
		}
	}
	log.Info("[server]: Redis closed")
}

func dbShutdown(ctx context.Context) {
	// close MongoDB
	if mgDBInstant != nil {
		if err := mgDBInstant.Disconnect(ctx); err != nil {
			log.Errorf("[MongoDB] shutdown error: %v", err)
		}
	}

	// close Postgres
	if pgDBInstant != nil {
		sqlDB, err := pgDBInstant.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Errorf("[Postgres] shutdown error: %v", err)
			}
		}
	}
	log.Info("[server]: DB closed")
}

func gracefulShutdown(ctx context.Context, server *http.Server) {
	log.Info("[server]: Shutting down server...")
	// close HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("HTTP server Shutdown: %v", err)
	}

	rabbitShutdown()
	redisShutdown()
	dbShutdown(ctx)
}
