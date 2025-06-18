package app_dependencies

import (
	"go-rebuild/internal/cache"
	"go-rebuild/internal/db"
	"go-rebuild/internal/mail"
	messagebroker "go-rebuild/internal/message_broker"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type CoreDependencies struct {
	DBRepo              db.DB
	RedisClient         *redis.Client // แก้ชื่อให้สื่อความหมาย
	CacheService        cache.Cache
	MailService         mail.Mail
	RabbitMQConn        *amqp.Connection
	ProducerChannel     *amqp.Channel
	UserConsumeChannel  *amqp.Channel
	StockConsumeChannel *amqp.Channel
	ProducerService     messagebroker.ProducerService
}


