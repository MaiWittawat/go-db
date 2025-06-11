package main

import (
	"context"
	appcore_config "go-rebuild/cmd/go-rebuild/config"
	"go-rebuild/internal/auth"
	rclient "go-rebuild/internal/cache"
	"go-rebuild/internal/db"
	"go-rebuild/internal/handler"
	"go-rebuild/internal/handler/api"
	"go-rebuild/internal/mail"
	messagebroker "go-rebuild/internal/message_broker"

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
	rabbitMQChannel    *amqp.Channel
)

func main() {
	overallStart := time.Now()
	start := time.Now()

	appcore_config.InitConfigurations()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	gin.SetMode(gin.ReleaseMode)

	if appcore_config.Config.Mode == "develop" {
		log.SetLevel(log.InfoLevel)
	}else {
		log.SetLevel(log.WarnLevel)
	}

	var dbRepo db.DB
	var err error
	useMongo := false

	if useMongo {
		start = time.Now()
		initMongoCtx, initMongoCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer initMongoCancel()

		mgDBInstant, err = db.InitMongoDB(initMongoCtx)
		if err != nil {
			log.Panic("fail to connect mongodb: ", err)
		}

		log.Printf("Startup: MongoDB connection took %s", time.Since(start))
		dbRepo = db.NewMongoRepo(mgDBInstant, "miniproject")

	} else {
		start = time.Now()
		pgDBInstant, err = db.InitPsqlDB()

		if err != nil {
			log.Panic("fail to connect psqldb: ", err)
		}

		log.Printf("Startup: PostgreSQL connection took %s", time.Since(start))
		dbRepo = db.NewPsqlRepo(pgDBInstant)
	}

	// init redis client
	start = time.Now()
	redisClientInstant = rclient.InitRedisClient()
	log.Printf("Startup: Redis client initialization took %s", time.Since(start))
	cacheSvc := rclient.NewCacheService(redisClientInstant)
	router := gin.Default()

	// init mail client
	start = time.Now()
	mailClient := mail.InitMailClient()
	log.Printf("Startup: Mail client initialization took %s", time.Since(start))
	mailService := mail.NewMailService(mailClient)

	// init rabbitmq
	start = time.Now()
	rabbitMQConn = messagebroker.InitRabbitmq()
	log.Printf("Startup: RabbitMQ connection took %s", time.Since(start))

	start = time.Now()
	rabbitMQChannel = messagebroker.OpenChannel(rabbitMQConn)
	log.Printf("Startup: RabbitMQ channel open took %s", time.Since(start))

	start = time.Now()
	userCreateQueueSetup(rabbitMQChannel)
	userUpdateQueueSetup(rabbitMQChannel)
	log.Printf("Startup: RabbitMQ queue setup took %s", time.Since(start))

	start = time.Now()
	producerService := messagebroker.NewProducerService(rabbitMQChannel)
	consumeUserService := messagebroker.NewConsumerService(rabbitMQConn, rabbitMQChannel, mailService)
	log.Printf("Startup: Message broker services creation took %s", time.Since(start))

	// User and Auth
	start = time.Now()
	userRepository := userRepo.NewUserRepo(dbRepo, cacheSvc)
	authService := auth.NewAuthService(userRepository, producerService)
	userService := userSvc.NewUserService(userRepository, producerService)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	api.RegisterUserAPI(router, userHandler, authService)
	api.RegisterAuthAPI(router, authHandler)
	log.Printf("Startup: User and Auth services setup took %s", time.Since(start))

	// Product
	start = time.Now()
	productRepo := productRepo.NewProductRepo(dbRepo, cacheSvc)
	productSvc := productSvc.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)
	api.RegisterProductAPI(router, productHandler, authService)
	log.Printf("Startup: Product services setup took %s", time.Since(start))

	start = time.Now()
	// Order
	orderRepo := orderRepo.NewOrderRepo(dbRepo, cacheSvc)
	orderSvc := orderSvc.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderSvc)
	api.RegisterOrderAPI(router, orderHandler, authService)
	log.Printf("Startup: Order services setup took %s", time.Since(start))

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

	go consumeUserService.Consuming(userSvc.QueueName, "user_consume")

	log.Printf("Overall application startup took %s", time.Since(overallStart))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown signal received")

	gracefulShutdown(shutdownCtx, server)

}

func userCreateQueueSetup(ch *amqp.Channel) {
	messagebroker.DeclareExchange(ch, userSvc.ExchangeName, userSvc.ExchangeType)
	messagebroker.DeclareQueue(ch, userSvc.QueueName)
	messagebroker.BindQueueToExchange(ch, userSvc.QueueName, userSvc.ExchangeName, "user.create")
}

func userUpdateQueueSetup(ch *amqp.Channel) {
	messagebroker.DeclareExchange(ch, userSvc.ExchangeName, userSvc.ExchangeType)
	messagebroker.DeclareQueue(ch, userSvc.QueueName)
	messagebroker.BindQueueToExchange(ch, userSvc.QueueName, userSvc.ExchangeName, "user.update")
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

	// close rabbitmq
	if rabbitMQChannel != nil {
		if err := rabbitMQChannel.Close(); err != nil {
			log.Errorf("RabbitMQ channel shutdown error: %v", err)
		} else {
			log.Info("RabbitMQ channel closed")
		}
	}
	if rabbitMQConn != nil {
		if err := rabbitMQConn.Close(); err != nil {
			log.Errorf("RabbitMQ connection shutdown error: %v", err)
		} else {
			log.Info("RabbitMQ connection closed")
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
