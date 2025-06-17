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

type consumerService struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	mailSvc  mail.Mail
	stockSvc module.StockService
}

func NewConsumerService(conn *amqp.Connection, ch *amqp.Channel, mailSvc mail.Mail, stockSvc module.StockService) ConsumerService {
	return &consumerService{
		conn:     conn,
		ch:       ch,
		mailSvc:  mailSvc,
		stockSvc: stockSvc,
	}
}

func (s *consumerService) EmailConsuming(queueName string, tag string) error {
	log.Printf("[consume]: %s called", tag)
	msgs, err := s.ch.Consume(
		queueName,
		tag,
		true,  // autoAck
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

			log.Println("user: ", user)

			switch msg.RoutingKey {
			case "user.create":
				email := []string{string(user.Email)}
				if err = s.mailSvc.SendWelcomeEmail(email); err != nil {
					log.WithError(err).Error("user consume created fail")
					continue
				}
				log.Printf("[Consume]: Received by Consumer '%s': user created\n", msg.ConsumerTag)

			case "user.update":
				email := []string{string(user.Email)}
				subject := "User Update"
				message := fmt.Sprintf("Your account %s has updated in go-rebuild project At %v", user.Email, user.UpdatedAt)
				if err = s.mailSvc.SendEmail(message, subject, email); err != nil {
					log.WithError(err).Error("user consume updated fail")
					continue
				}
				log.Printf("[Consume]: Received by Consumer '%s': user updated\n", msg.ConsumerTag)

			default:
				log.Printf("[Consume]: unsupported message type: %s\n", msg.RoutingKey)
			}
		}
	}()

	return nil
}

func (s *consumerService) StockConsuming(queueName string, tag string) error {
	log.Printf("[Consume]: %s called", tag)
	msgs, err := s.ch.Consume(
		queueName,
		tag,
		true,  // autoAck
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
			log.Println("msg: ", msg)
			var stock model.Stock
			if err := json.Unmarshal(msg.Body, &stock); err != nil {
				log.WithError(err).Error("fail to unmarshal stock")
				continue
			}

			switch msg.RoutingKey {
			case "stock.create":
				if err := s.stockSvc.Save(context.Background(), stock.ProductID, stock.Quantity); err != nil {
					log.WithError(err).Error("stock consume save failed")
					continue
				}
				log.Printf("[Consume]: Stock created received by consumer '%s': stock create", msg.ConsumerTag)

			case "stock.update":
				log.Printf("[Consume in update stock]: Received by Consumer '%s': stock updated", msg.ConsumerTag)
				if err := s.stockSvc.Update(context.Background(), stock.ProductID, stock.Quantity); err != nil {
					log.WithError(err).Error("stock consume update failed")
					continue
				}
				log.Printf("[Consume]: Received by Consumer '%s': stock updated", msg.ConsumerTag)

			case "stock.increase":
				if err := s.stockSvc.IncreaseQuantity(context.Background(), stock.Quantity, stock.ProductID); err != nil {
					log.WithError(err).Error("stock consume increase failed")
					continue
				}
				log.Printf("[Consume]: Received by Consumer '%s': stock increase", msg.ConsumerTag)

			case "stock.decrease":
				if err := s.stockSvc.DecreaseQuantity(context.Background(), stock.Quantity, stock.ProductID); err != nil {
					log.WithError(err).Error("stock consume decrease failed")
					continue
				}
				log.Printf("[Consume]: Received by Consumer '%s': stock decrease", msg.ConsumerTag)

			default:
				log.Printf("[Consume]: Unsupported message type: %s", msg.RoutingKey)
			}
		}
	}()

	return nil
}
