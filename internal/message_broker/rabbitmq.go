package messagebroker

import (
	"context"
	appcore_config "go-rebuild/cmd/go-rebuild/config"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type MQConfig struct {
	ExchangeName string
	ExchangeType string
	QueueName    string
	RoutingKey   string
}

type ConsumerService interface {
	Consuming(queueName string, tag string) error
}

type ProducerService interface {
	Publishing(ctx context.Context, mqConf *MQConfig, body []byte) error
}

func HandleError(err error, msg string) {
	if err != nil {
		log.Printf("%s:, %v", msg, err)
	}
}

func NewMQConfig(exName, exType, qName, routingKey string) *MQConfig {
	return &MQConfig{
		ExchangeName: exName,
		ExchangeType: exType,
		QueueName:    qName,
		RoutingKey:   routingKey,
	}
}

func InitRabbitmq() *amqp.Connection {
	conn, err := amqp.Dial(appcore_config.Config.RabbitmqUrl)
	HandleError(err, "fail to connect rabbitmq")
	return conn
}

func SetupExchangeAndQueue(ch *amqp.Channel, cfg *MQConfig) error {
	DeclareExchange(ch, cfg.ExchangeName, cfg.ExchangeType)
	DeclareQueue(ch, cfg.QueueName)
	BindQueueToExchange(ch, cfg.QueueName, cfg.ExchangeName, cfg.RoutingKey)
	return nil
}

func OpenChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	HandleError(err, "fail to opend channel")
	return ch
}

func DeclareQueue(ch *amqp.Channel, queueName string) *amqp.Queue {
	queue, err := ch.QueueDeclare(
		queueName,
		true, // durable
		false,
		false,
		false,
		nil,
	)

	HandleError(err, "failed to declare queue")
	return &queue
}

func DeclareExchange(ch *amqp.Channel, exchangeName, exchangeType string) {
	err := ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	HandleError(err, "Failed to declare an exchange")
	log.Printf("Exchange '%s' declared successfully\n", exchangeName)
}

func BindQueueToExchange(ch *amqp.Channel, queueName, exchangeName, routingKey string) {
	err := ch.QueueBind(
		queueName,    // name of the queue
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // noWait
		nil,          // arguments
	)
	HandleError(err, "Failed to bind a queue to an exchange")
	log.Printf("Queue '%s' bound to Exchange '%s' with RoutingKey '%s'\n", queueName, exchangeName, routingKey)
}
