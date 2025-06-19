package app_setup

import (
	"context"
	"fmt"
	appcore_config "go-rebuild/cmd/go-rebuild/config"
	"go-rebuild/internal/auth"
	redisclient "go-rebuild/internal/cache"
	dbRepo "go-rebuild/internal/db"
	"go-rebuild/internal/handler"
	"go-rebuild/internal/handler/api"
	authHandler "go-rebuild/internal/handler/auth"
	messageHandler "go-rebuild/internal/handler/message"
	orderHandler "go-rebuild/internal/handler/order"
	productHandler "go-rebuild/internal/handler/product"
	stockHandler "go-rebuild/internal/handler/stock"
	storerHandler "go-rebuild/internal/handler/storer"
	userHandler "go-rebuild/internal/handler/user"
	"go-rebuild/internal/mail"
	messagebroker "go-rebuild/internal/message_broker"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	messageSvc "go-rebuild/internal/module/message"
	orderSvc "go-rebuild/internal/module/order"
	productSvc "go-rebuild/internal/module/product"
	stockSvc "go-rebuild/internal/module/stock"
	userSvc "go-rebuild/internal/module/user"
	"go-rebuild/internal/realtime"
	"go-rebuild/internal/repository"
	messageRepo "go-rebuild/internal/repository/message"
	orderRepo "go-rebuild/internal/repository/order"
	productRepo "go-rebuild/internal/repository/product"
	stockRepo "go-rebuild/internal/repository/stock"
	userRepo "go-rebuild/internal/repository/user"
	"go-rebuild/internal/storer"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type AppClients struct {
	MongoDB      *mongo.Client
	PostgresDB   *gorm.DB
	RedisClient  *redis.Client
	RabbitMQConn *amqp.Connection

	// Auto close
	MailClient          *gomail.Dialer
	WebSocketServer     *realtime.WebSocketServer
	MinioClient         *minio.Client
	ProducerChannel     *amqp.Channel
	UserConsumeChannel  *amqp.Channel
	StockConsumeChannel *amqp.Channel
}

type AppServices struct {
	// Repositories
	DBRepository      dbRepo.DB
	UserRepository    repository.UserRepository
	ProductRepository repository.ProductRepository
	OrderRepository   repository.OrderRepository
	StockRepository   repository.StockRepository
	MessageRepository repository.MessageRepository

	// Services
	AuthService     auth.Jwt
	UserService     module.UserService
	StockService    module.StockService
	ProductService  module.ProductService
	ProducerService messagebroker.ProducerService
	ConsumerService messagebroker.ConsumerService
	MessageService  module.MessageService
	MailService     mail.Mail
	OrderService    module.OrderService
	StorageService  storer.Storer
	MQBroker        messagebroker.MessageBroker
	LiveChat        *realtime.LiveChat

	// Handlers
	AuthHandler    handler.AuthHandler
	UserHandler    handler.UserHandler
	ProductHandler handler.ProductHandler
	OrderHandler   handler.OrderHandler
	StockHandler   handler.StockHandler
	MessageHandler handler.MessageHandler
	StorerHandler  handler.StorerHandler
}

// ------------------------------ Init 3rd Party ------------------------------
func InitAppClients() (*AppClients, error) {
	clients := &AppClients{}
	var err error

	// init Postgres
	clients.PostgresDB, err = dbRepo.InitPsqlDB()
	if err != nil {
		return nil, err
	}

	// Init Redis
	clients.RedisClient = redisclient.InitRedisClient(appcore_config.Config.RedisUrl, appcore_config.Config.RedisPass)

	// Init Mail
	clients.MailClient = mail.InitSMTP()

	// Init WebSocket
	clients.WebSocketServer = realtime.NewWebSocketServer()

	// Init Minio
	clients.MinioClient = storer.InitMinio(appcore_config.Config.MinioURL, appcore_config.Config.MinioAccessKey, appcore_config.Config.MinioSecretKey)

	// Init RabbitMQ
	clients.RabbitMQConn = messagebroker.InitRabbitmq()
	clients.ProducerChannel = messagebroker.OpenChannel(clients.RabbitMQConn)
	clients.UserConsumeChannel = messagebroker.OpenChannel(clients.RabbitMQConn)
	clients.StockConsumeChannel = messagebroker.OpenChannel(clients.RabbitMQConn)

	// Setup RabbitMQ Queues
	if err := messagebroker.SetupExchangeAndQueue(clients.UserConsumeChannel,
		&model.MQConfig{
			ExchangeName: messagebroker.UserExchangeName,
			ExchangeType: messagebroker.UserExchangeType,
			QueueName:    messagebroker.UserQueueName,
			RoutingKey:   "user.#",
		}); err != nil {
		return nil, fmt.Errorf("failed to setup user exchange and queue: %w", err)
	}

	if err := messagebroker.SetupExchangeAndQueue(clients.StockConsumeChannel,
		&model.MQConfig{
			ExchangeName: messagebroker.StockExchangeName,
			ExchangeType: messagebroker.StockExchangeType,
			QueueName:    messagebroker.StockQueueName,
			RoutingKey:   "stock.#",
		}); err != nil {
		return nil, fmt.Errorf("failed to setup stock exchange and queue: %w", err)
	}

	return clients, nil
}

// ------------------------------ Build Application Service ------------------------------
func BuildApplicationServices(clients *AppClients) *AppServices {
	cacheService := redisclient.NewCacheService(clients.RedisClient)
	dbRepository := dbRepo.NewPsqlRepo(clients.PostgresDB)

	// Repository
	userRepository := userRepo.NewUserRepo(dbRepository, cacheService)
	productRepository := productRepo.NewProductRepo(dbRepository, cacheService)
	orderRepository := orderRepo.NewOrderRepo(dbRepository, cacheService)
	stockRepository := stockRepo.NewStockRepo(dbRepository, cacheService)
	messageRepository := messageRepo.NewMessageRepo(dbRepository, cacheService)

	// Service
	mailService := mail.NewMailService(clients.MailClient)
	stockService := stockSvc.NewStockService(stockRepository)
	producerService := messagebroker.NewProducer(clients.ProducerChannel)
	consumerService := messagebroker.NewConsumer(clients.UserConsumeChannel, clients.StockConsumeChannel, mailService, stockService)
	mqBroker := messagebroker.NewMessageBroker(producerService, consumerService)
	userService := userSvc.NewUserService(userRepository, producerService)
	authService := auth.NewAuthService(userService, producerService)
	productService := productSvc.NewProductService(productRepository, producerService)
	orderService := orderSvc.NewOrderService(orderRepository, productService, producerService)
	messageService := messageSvc.NewMessageService(messageRepository)
	liveChat := realtime.NewLiveChat(clients.WebSocketServer, messageService, authService)
	storageService := storer.NewStorerService(clients.MinioClient, appcore_config.Config.MinioBucketName)

	// Handler
	authHandler := authHandler.NewAuthHandler(authService)
	userHandler := userHandler.NewUserHandler(userService)
	productHandler := productHandler.NewProductHandler(productService)
	orderHandler := orderHandler.NewOrderHandler(orderService)
	stockHandler := stockHandler.NewStockHandler(stockService)
	messageHandler := messageHandler.NewMessageHandler(liveChat, messageService)
	storerHandler := storerHandler.NewStorerHandler(storageService)

	return &AppServices{
		DBRepository:      dbRepository,
		UserRepository:    userRepository,
		ProductRepository: productRepository,
		OrderRepository:   orderRepository,
		StockRepository:   stockRepository,
		MessageRepository: messageRepository,
		StockService:      stockService,
		ProducerService:   producerService,
		ConsumerService:   consumerService,
		MailService:       mailService,
		MQBroker:          mqBroker,
		UserService:       userService,
		AuthService:       authService,
		ProductService:    productService,
		OrderService:      orderService,
		MessageService:    messageService,
		LiveChat:          liveChat,
		StorageService:    storageService,
		AuthHandler:       authHandler,
		UserHandler:       userHandler,
		ProductHandler:    productHandler,
		OrderHandler:      orderHandler,
		StockHandler:      stockHandler,
		MessageHandler:    messageHandler,
		StorerHandler:     storerHandler,
	}
}

// ------------------------------ API Routing ------------------------------
func APIRoutes(router *gin.Engine, appSvc *AppServices) {
	apiConf := api.NewAPIRouterConfigurator(
		appSvc.AuthHandler,
		appSvc.UserHandler,
		appSvc.OrderHandler,
		appSvc.ProductHandler,
		appSvc.StockHandler,
		appSvc.MessageHandler,
		appSvc.StorerHandler,
		appSvc.AuthService,
	)

	apiConf.PublicAPIRoutes(router)
	apiConf.ProtectAPIRoutes(router)
}

// ------------------------------ Message Queue Consumers ------------------------------
func StartConsumers(broker messagebroker.MessageBroker) {
	go broker.EmailConsuming(messagebroker.UserQueueName, "user_consume")
	go broker.StockConsuming(messagebroker.StockQueueName, "stock_consume")
}

// ------------------------------ Graceful Shutdown Functions ------------------------------
func GracefulShutdown(ctx context.Context, server *http.Server, clients *AppClients) {
	log.Info("[server]: Shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("[server]: HTTP server Shutdown %v", err)
	}

	rabbitShutdown(clients.RabbitMQConn)
	redisShutdown(clients.RedisClient)
	dbShutdown(ctx, clients.MongoDB, clients.PostgresDB)
}

func rabbitShutdown(rabbitMQConn *amqp.Connection) {
	// close rabbitmq connection
	if rabbitMQConn != nil {
		if err := rabbitMQConn.Close(); err != nil {
			log.Errorf("[RabbitMQ] connection shutdown error: %v", err)
		}
	}
	log.Info("[server]: Rabbitmq closed")
}

func redisShutdown(redisConn *redis.Client) {
	// close Redis
	if redisConn != nil {
		if err := redisConn.Close(); err != nil {
			log.Errorf("[Redis] shutdown error: %v", err)
		}
	}
	log.Info("[server]: Redis closed")
}

func dbShutdown(ctx context.Context, mgClient *mongo.Client, psqlDB *gorm.DB) {
	// close MongoDB
	if mgClient != nil {
		if err := mgClient.Disconnect(ctx); err != nil {
			log.Errorf("[MongoDB] shutdown error: %v", err)
		}
	}

	// close Postgres
	if psqlDB != nil {
		sqlDB, err := psqlDB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Errorf("[Postgres] shutdown error: %v", err)
			}
		}
	}
	log.Info("[server]: DB closed")
}
