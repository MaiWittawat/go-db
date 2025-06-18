package messagebroker

import (
	"context"
	"encoding/json"
	"fmt"
	"go-rebuild/internal/mail"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type messageBroker struct {
	producerService ProducerService
	consumerService ConsumerService
}

type consumerService struct {
	mailSvc  mail.Mail
	stockSvc module.StockService
	userCh   *amqp.Channel
	stockCh  *amqp.Channel
}

type producerService struct {
	ch *amqp.Channel
}

// ------------------------ Message Broker ------------------------
func NewMessageBroker(producerSvc ProducerService, consumerSvc ConsumerService) MessageBroker {
	return &messageBroker{
		producerService: producerSvc,
		consumerService: consumerSvc,
	}
}

func (m *messageBroker) Publishing(ctx context.Context, mqConf *model.MQConfig, body []byte) error {
	return m.producerService.Publishing(ctx, mqConf, body)
}

func (m *messageBroker) EmailConsuming(queueName string, tag string) error {
	return m.consumerService.EmailConsuming(queueName, tag)
}

func (m *messageBroker) StockConsuming(queueName string, tag string) error {
	return m.consumerService.StockConsuming(queueName, tag)
}

// ------------------------ Publisher ------------------------
func NewProducer(ch *amqp.Channel) ProducerService {
	return &producerService{
		ch: ch,
	}
}

func (p *producerService) Publishing(ctx context.Context, mqConf *model.MQConfig, body []byte) error {
	if err := p.ch.PublishWithContext(
		ctx,
		mqConf.ExchangeName,
		mqConf.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: 2, // persistant
		},
	); err != nil {
		return err
	}
	log.Info("[Publisher]: Publish success")
	return nil
}

// ------------------------ Consumer ------------------------
func NewConsumer(userCh *amqp.Channel, stockCh *amqp.Channel, mailSvc mail.Mail, stockSvc module.StockService) ConsumerService {
	return &consumerService{
		mailSvc:  mailSvc,
		stockSvc: stockSvc,
		userCh:   userCh,
		stockCh:  stockCh,
	}
}

func (c *consumerService) EmailConsuming(queueName string, tag string) error {
	log.Printf("[consume]: %s called", tag)
	msgs, err := c.userCh.Consume(
		queueName,
		tag,
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var user model.User
			if err := json.Unmarshal(msg.Body, &user); err != nil {
				log.WithError(err).Error("fail to unmarshal user")
				continue
			}

			log.Println("user email: ", user.Email)

			switch msg.RoutingKey {
			case "user.create":
				email := []string{string(user.Email)}
				if err = c.mailSvc.SendWelcomeEmail(email); err != nil {
					log.WithError(err).Error("user consume created fail")
					continue
				}
				msg.Ack(false)
				log.Printf("[Consume]: Received by Consumer '%s': user created\n", msg.ConsumerTag)

			case "user.update":
				email := []string{string(user.Email)}
				subject := "User Update"
				message := fmt.Sprintf("Your account %s has updated in go-rebuild project At %v", user.Email, user.UpdatedAt)
				if err = c.mailSvc.SendEmail(message, subject, email); err != nil {
					log.WithError(err).Error("user consume updated fail")
					continue
				}

				msg.Ack(false)
				log.Printf("[Consume]: Received by Consumer '%s': user updated\n", msg.ConsumerTag)

			default:
				log.Printf("[Consume]: unsupported message type: %s\n", msg.RoutingKey)
			}
		}
	}()

	return nil
}

func (c *consumerService) StockConsuming(queueName string, tag string) error {
	log.Printf("[Consume]: %s called", tag)
	msgs, err := c.stockCh.Consume(
		queueName,
		tag,
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var stock model.Stock
			if err := json.Unmarshal(msg.Body, &stock); err != nil {
				log.WithError(err).Error("fail to unmarshal stock")
				continue
			}

			switch msg.RoutingKey {
			case "stock.create":
				if err := c.stockSvc.Save(context.Background(), stock.ProductID, stock.Quantity); err != nil {
					log.WithError(err).Error("stock consume save failed")
					continue
				}
				msg.Ack(false)
				log.Printf("[Consume]: Stock created received by consumer '%s': stock create", msg.ConsumerTag)

			case "stock.update":
				log.Printf("[Consume in update stock]: Received by Consumer '%s': stock updated", msg.ConsumerTag)
				if err := c.stockSvc.Update(context.Background(), stock.ProductID, stock.Quantity); err != nil {
					log.WithError(err).Error("stock consume update failed")
					continue
				}
				msg.Ack(false)
				log.Printf("[Consume]: Received by Consumer '%s': stock updated", msg.ConsumerTag)

			case "stock.increase":
				if err := c.stockSvc.IncreaseQuantity(context.Background(), stock.Quantity, stock.ProductID); err != nil {
					log.WithError(err).Error("stock consume increase failed")
					continue
				}
				msg.Ack(false)
				log.Printf("[Consume]: Received by Consumer '%s': stock increase", msg.ConsumerTag)

			case "stock.decrease":
				if err := c.stockSvc.DecreaseQuantity(context.Background(), stock.Quantity, stock.ProductID); err != nil {
					log.WithError(err).Error("stock consume decrease failed")
					continue
				}
				msg.Ack(false)
				log.Printf("[Consume]: Received by Consumer '%s': stock decrease", msg.ConsumerTag)

			default:
				log.Printf("[Consume]: Unsupported message type: %s", msg.RoutingKey)
			}
		}
	}()

	return nil
}
