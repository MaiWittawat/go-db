package messagebroker

import (
	"context"
	"go-rebuild/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type producerService struct {
	ch *amqp.Channel
}

func NewProducerService(ch *amqp.Channel) ProducerService {
	return &producerService{ch: ch}
}

func (s *producerService) Publishing(ctx context.Context, mqConf *model.MQConfig, body []byte) error {
	if err := s.ch.PublishWithContext(
		ctx,
		mqConf.ExchangeName,
		mqConf.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	); err != nil {
		return err
	}
	log.Info("[Publisher]: Publish success")
	return nil
}