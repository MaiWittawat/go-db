package messagebroker

import (
	"context"
	appcore_config "go-rebuild/cmd/go-rebuild/config"
	"go-rebuild/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)


var (
	// user queue
	UserExchangeName = "user_exchange"
	UserExchangeType = "topic"
	UserQueueName    = "user_queue"

	// stock queue
	StockExchangeName = "stock_exchange"
	StockExchangeType = "topic"
	StockQueueName    = "stock_queue"
)


type MessageBroker interface {
	ConsumerService
	ProducerService
}

type ConsumerService interface {
	EmailConsuming(queueName string, tag string) error
	StockConsuming(queueName string, tag string) error
}

type ProducerService interface {
	Publishing(ctx context.Context, mqConf *model.MQConfig, body []byte) error
}

func HandleError(err error, msg string) {
	if err != nil {
		log.Printf("[Error]: %s, %v", msg, err)
	}
}

func InitRabbitmq() *amqp.Connection {
	conn, err := amqp.Dial(appcore_config.Config.RabbitmqUrl)
	HandleError(err, "fail to connect rabbitmq")
	return conn
}

func SetupExchangeAndQueue(ch *amqp.Channel, cfg *model.MQConfig) error {
	if err := DeclareExchange(ch, cfg.ExchangeName, cfg.ExchangeType); err != nil {return err}
	if err := DeclareQueue(ch, cfg.QueueName); err != nil {return err}
	if err := BindQueueToExchange(ch, cfg.QueueName, cfg.ExchangeName, cfg.RoutingKey); err != nil {return err}
	return nil
}

func OpenChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	HandleError(err, "fail to opend channel")
	return ch
}

func DeclareQueue(ch *amqp.Channel, queueName string) error {
	_, err := ch.QueueDeclare(
		queueName,
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeclareExchange(ch *amqp.Channel, exchangeName, exchangeType string) error {
	err := ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

func BindQueueToExchange(ch *amqp.Channel, queueName, exchangeName, routingKey string) error {
	err := ch.QueueBind(
		queueName,    // name of the queue
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		return err
	}
	return nil
}
