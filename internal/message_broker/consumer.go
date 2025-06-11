package messagebroker

import (
	"encoding/json"
	"fmt"
	"go-rebuild/internal/mail"
	"go-rebuild/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type consumerService struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	mailSvc mail.Mail
}

func NewConsumerService(conn *amqp.Connection, ch *amqp.Channel, mailSvc mail.Mail) ConsumerService {
	return &consumerService{conn: conn, ch: ch, mailSvc: mailSvc}
}

func (s *consumerService) Consuming(queueName string, tag string) error {
	log.Printf("consume: %s called", tag)
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
			var envelope model.EnvelopeBroker
			if err := json.Unmarshal(msg.Body, &envelope); err != nil {
				log.Printf("invalid envelope: %v", err)
				continue
			}

			switch envelope.Type {
			// case "create_order":
			// 	var order model.Order
			// 	if err := json.Unmarshal(envelope.Payload, &order); err != nil {
			// 		log.WithError(err).Error("fail to unmarshal order")
			// 		continue
			// 	}

			// 	subject := "Order Created"
			// 	email := []string{string(order.Email)}
			// 	message := fmt.Sprintf("hello you are buy product: %s, quantity: %d, at %v", order.ProductName, order.Quantity, order.CreatedAt)
			// 	s.mailSvc.SendEmail(message, subject, email)
			// 	log.Printf("Received by Consumer '%s': %s", msg.ConsumerTag, msg.Body)

			case "create_user":
				var user model.User
				if err := json.Unmarshal(envelope.Payload, &user); err != nil {
					log.WithError(err).Error("fail to unmarshal user")
					continue
				}

				email := []string{string(user.Email)}
				s.mailSvc.SendWelcomeEmail(email)
				log.Printf("Received by Consumer '%s': user created", msg.ConsumerTag)

			case "update_user":
				var user model.User
				if err := json.Unmarshal(envelope.Payload, &user); err != nil {
					log.WithError(err).Error("fail to unmarshal user")
					continue
				}

				email := []string{string(user.Email)}
				subject := "User Update"
				message := fmt.Sprintf("Your account %s has updated in go-rebuild project At %v", user.Email, user.UpdatedAt)
				s.mailSvc.SendEmail(message, subject, email)
				log.Printf("Received by Consumer '%s': user updated", msg.ConsumerTag)

			default:
				log.Printf("unsupported message type: %s", envelope.Type)
			}
		}
	}()

	return nil
}
