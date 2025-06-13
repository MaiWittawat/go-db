package messagebroker

import (
	"encoding/json"
	"fmt"
	"go-rebuild/internal/mail"
	"go-rebuild/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type emailConsumerService struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	mailSvc mail.Mail
}

func NewEmailConsumerService(conn *amqp.Connection, ch *amqp.Channel, mailSvc mail.Mail) ConsumerService {
	return &emailConsumerService{conn: conn, ch: ch, mailSvc: mailSvc}
}

func (s *emailConsumerService) Consuming(queueName string, tag string) error {
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
			var envelope model.Envelope
			if err := json.Unmarshal(msg.Body, &envelope); err != nil {
				log.Printf("[Consume]: invalid envelope: %v\n", err)
				continue
			}
			var user model.User
			if err := json.Unmarshal(envelope.Payload, &user); err != nil {
				log.WithError(err).Error("fail to unmarshal user")
				continue
			}

			switch envelope.Type {
			case "create_user":
				email := []string{string(user.Email)}
				if err = s.mailSvc.SendWelcomeEmail(email); err != nil {
					log.WithError(err).Error("user consume created fail")
					continue
				}
				log.Printf("[Consume]: Received by Consumer '%s': user created\n", msg.ConsumerTag)

			case "update_user":
				email := []string{string(user.Email)}
				subject := "User Update"
				message := fmt.Sprintf("Your account %s has updated in go-rebuild project At %v", user.Email, user.UpdatedAt)
				if err = s.mailSvc.SendEmail(message, subject, email); err != nil {
					log.WithError(err).Error("user consume updated fail")
					continue
				}
				log.Printf("[Consume]: Received by Consumer '%s': user updated\n", msg.ConsumerTag)

			default:
				log.Printf("[Consume]: unsupported message type: %s\n", envelope.Type)
			}
		}
	}()

	return nil
}
