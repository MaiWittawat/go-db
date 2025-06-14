package messagebroker

import (
	"context"
	"encoding/json"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type stockComsumeService struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	stockSvc module.StockService
}

func NewStockComsumeService(conn *amqp.Connection, ch *amqp.Channel, stockSvc module.StockService) ConsumerService {
	return &stockComsumeService{conn: conn, ch: ch, stockSvc: stockSvc}
}

func (s *stockComsumeService) Consuming(queueName string, tag string) error {
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
			var envelope model.Envelope
			if err := json.Unmarshal(msg.Body, &envelope); err != nil {
				log.Printf("[Consume]: invalid envelope: %v", err)
				continue
			}

			var stock model.Stock
			if err := json.Unmarshal(envelope.Payload, &stock); err != nil {
				log.WithError(err).Error("fail to unmarshal stock")
				continue
			}

			switch envelope.Type {
			case "create_stock":
				if err := s.stockSvc.Save(context.Background(), stock.ProductID, stock.Quantity); err != nil {
					log.WithError(err).Error("stock consume save failed")
					continue
				}
				log.Printf("[Consume]: Stock created received by consumer '%s': stock create", msg.ConsumerTag)

			case "update_stock":
				log.Printf("[Consume in update stock]: Received by Consumer '%s': stock updated", msg.ConsumerTag)
				if err := s.stockSvc.Update(context.Background(), stock.ProductID, stock.Quantity); err != nil {
					log.WithError(err).Error("stock consume update failed")
					continue
				}
				log.Printf("[Consume]: Received by Consumer '%s': stock updated", msg.ConsumerTag)

			case "increase_stock":
				if err := s.stockSvc.IncreaseQuantity(context.Background(), stock.Quantity, stock.ProductID); err != nil {
					log.WithError(err).Error("stock consume increase failed")
					continue
				}
				log.Printf("[Consume]: Received by Consumer '%s': stock increase", msg.ConsumerTag)

			case "decrease_stock":
				if err := s.stockSvc.DecreaseQuantity(context.Background(), stock.Quantity, stock.ProductID); err != nil {
					log.WithError(err).Error("stock consume decrease failed")
					continue
				}
				log.Printf("[Consume]: Received by Consumer '%s': stock decrease", msg.ConsumerTag)

			default:
				log.Printf("[Consume]: Unsupported message type: %s", envelope.Type)
			}
		}
	}()

	return nil
}
