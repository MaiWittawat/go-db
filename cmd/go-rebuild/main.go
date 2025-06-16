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
	"go-rebuild/internal/realtime"

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
	mgDBInstant         *mongo.Client
	pgDBInstant         *gorm.DB
	redisClientInstant  *redis.Client
	rabbitMQConn        *amqp.Connection
	producerChannel     *amqp.Channel
	userConsumeChannel  *amqp.Channel
	stockConsumeChannel *amqp.Channel
)

func main() {
	// ------------------------------ Setup Config ------------------------------
	overallStart := time.Now()
	start := time.Now()

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
		start = time.Now()
		initMongoCtx, initMongoCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer initMongoCancel()

		mgDBInstant, err = db.InitMongoDB(initMongoCtx)
		if err != nil {
			log.Panic("fail to connect mongodb: ", err)
		}

		log.Printf("[Startup]: MongoDB connection took %s", time.Since(start))
		dbRepo = db.NewMongoRepo(mgDBInstant, "miniproject")

	} else {
		start = time.Now()
		pgDBInstant, err = db.InitPsqlDB()

		if err != nil {
			log.Panic("fail to connect psqldb: ", err)
		}

		log.Printf("[Startup]: PostgreSQL connection took %s", time.Since(start))
		dbRepo = db.NewPsqlRepo(pgDBInstant)
	}

	// ------------------------------ Init 3rd party ------------------------------
	// init redis client
	start = time.Now()
	redisClientInstant = rclient.InitRedisClient()
	log.Printf("[Startup]: Redis client initialization took %s", time.Since(start))
	cacheSvc := rclient.NewCacheService(redisClientInstant)
	router := gin.Default()

	// init mail client
	start = time.Now()
	mailClient := mail.InitMailClient()
	log.Printf("[Startup]: Mail client initialization took %s", time.Since(start))
	mailService := mail.NewMailService(mailClient)

	// init rabbitmq
	start = time.Now()

	// open 1 connection
	rabbitMQConn = messagebroker.InitRabbitmq()
	log.Printf("[Startup]: RabbitMQ connection took %s", time.Since(start))

	// open channel depen on service(use-case)
	start = time.Now()

	// rabbitMQChannel = messagebroker.OpenChannel(rabbitMQConn)
	producerChannel = messagebroker.OpenChannel(rabbitMQConn)
	userConsumeChannel = messagebroker.OpenChannel(rabbitMQConn)
	stockConsumeChannel = messagebroker.OpenChannel(rabbitMQConn)

	log.Printf("[Startup]: RabbitMQ channel open took %s", time.Since(start))

	start = time.Now()
	// user set rabbitmq up
	userCreateQueueSetup(userConsumeChannel)
	userUpdateQueueSetup(userConsumeChannel)
	// stock set rabbitmq up
	stockCreateQueueSetup(stockConsumeChannel)
	stockUpdateQueueSetup(stockConsumeChannel)
	log.Printf("[Startup]: RabbitMQ queue setup took %s", time.Since(start))

	start = time.Now()
	producerService := messagebroker.NewProducerService(producerChannel)
	log.Printf("[Startup]: Message broker services creation took %s", time.Since(start))

	// ------------------------------ Start service ------------------------------
	// User and Auth
	start = time.Now()
	userRepository := userRepo.NewUserRepo(dbRepo, cacheSvc)
	authService := auth.NewAuthService(userRepository, producerService)
	userService := userSvc.NewUserService(userRepository, producerService)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	api.RegisterUserAPI(router, userHandler, authService)
	api.RegisterAuthAPI(router, authHandler)
	log.Printf("[Startup]: User and Auth services setup took %s", time.Since(start))

	// Product && Stock
	// stock
	start = time.Now()
	stockRepository := stockRepo.NewStockRepo(dbRepo, cacheSvc)
	stockService := stockSvc.NewStockService(stockRepository)
	log.Printf("[Startup]: Stock services setup took %s", time.Since(start))

	// product
	start = time.Now()
	productRepo := productRepo.NewProductRepo(dbRepo, cacheSvc)
	productSvc := productSvc.NewProductService(productRepo, producerService)
	productHandler := handler.NewProductHandler(productSvc)
	api.RegisterProductAPI(router, productHandler, authService)
	log.Printf("[Startup]: Product services setup took %s", time.Since(start))

	// Order
	start = time.Now()
	orderRepository := orderRepo.NewOrderRepo(dbRepo, cacheSvc)
	orderService := orderSvc.NewOrderService(orderRepository, productSvc, producerService)
	orderHandler := handler.NewOrderHandler(orderService)
	api.RegisterOrderAPI(router, orderHandler, authService)
	log.Printf("[Startup]: Order services setup took %s", time.Since(start))

	// Message
	start = time.Now()
	messageRepository := messageRepo.NewMessageRepo(dbRepo, cacheSvc)
	messageService := messageSvc.NewMessageService(messageRepository)

	// Realtime service และ adapter
	websocketServer := realtime.NewWebSocketServer()
	chatRealtime := realtime.NewChatRealtime(websocketServer, messageService, authService)
	messageHandler := handler.NewMessageHandler(chatRealtime)
	api.RegisterMessageAPI(router, messageHandler, authService)
	log.Printf("[Startup]: Message services setup took %s", time.Since(start))

	userConsumeService := messagebroker.NewEmailConsumerService(rabbitMQConn, userConsumeChannel, mailService)
	stockConsumeService := messagebroker.NewStockComsumeService(rabbitMQConn, stockConsumeChannel, stockService)

	// ------------------------------ Start server ------------------------------
	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Info("[Down]: Server closed gracefully")
			} else {
				log.Fatalf("listen: %s\n", err)
			}
		}
	}()

	go userConsumeService.Consuming(userSvc.QueueName, "user_consume")
	go stockConsumeService.Consuming(stockSvc.QueueName, "stock_consume")

	log.Printf("[Time]: Overall application [startup] took %s", time.Since(overallStart))

	// ------------------------------ Shutdown ------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("[Signal]: shutdown signal received")

	// call shutdown all service
	gracefulShutdown(shutdownCtx, server)

}

// ------------------------------ Rabbitmq setup channel function ------------------------------
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

func stockCreateQueueSetup(ch *amqp.Channel) {
	messagebroker.DeclareExchange(ch, stockSvc.ExchangeName, stockSvc.ExchangeType)
	messagebroker.DeclareQueue(ch, stockSvc.QueueName)
	messagebroker.BindQueueToExchange(ch, stockSvc.QueueName, stockSvc.ExchangeName, "stock.create")
}

func stockUpdateQueueSetup(ch *amqp.Channel) {
	messagebroker.DeclareExchange(ch, stockSvc.ExchangeName, stockSvc.ExchangeType)
	messagebroker.DeclareQueue(ch, stockSvc.QueueName)
	messagebroker.BindQueueToExchange(ch, stockSvc.QueueName, stockSvc.ExchangeName, "stock.update")
}

// ------------------------------ Shutdown function ------------------------------
func rabbitShutdown() {
	// close rabbitmq producer consume channel
	if producerChannel != nil {
		if err := producerChannel.Close(); err != nil {
			log.Errorf("[RabbitMQ] producer consume channel shutdown error: %v", err)
		}
	}

	// close rabbitmq user consume channel
	if userConsumeChannel != nil {
		if err := userConsumeChannel.Close(); err != nil {
			log.Errorf("[RabbitMQ] user consume channel shutdown error: %v", err)
		}
	}

	// close rabbitmq user consume channel
	if stockConsumeChannel != nil {
		if err := stockConsumeChannel.Close(); err != nil {
			log.Errorf("[RabbitMQ] stock consume channel shutdown error: %v", err)
		}
	}

	// close rabbitmq connection
	if rabbitMQConn != nil {
		if err := rabbitMQConn.Close(); err != nil {
			log.Errorf("[RabbitMQ] connection shutdown error: %v", err)
		}
	}

	log.Info("[Down]: Rabbitmq closed")
}

func redisShutdown() {
	// close Redis
	if redisClientInstant != nil {
		if err := redisClientInstant.Close(); err != nil {
			log.Errorf("[Redis] shutdown error: %v", err)
		}
	}

	log.Info("[Down]: Redis closed")
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
	log.Info("[Down]: DB closed")
}

func gracefulShutdown(ctx context.Context, server *http.Server) {
	log.Info("[Down]: Shutting down server...")
	// close HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("HTTP server Shutdown: %v", err)
	}

	rabbitShutdown()
	redisShutdown()
	dbShutdown(ctx)
}
