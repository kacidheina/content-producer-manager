package rabbitmq

import (
	"content-producer-manager/pkg/model"
	"encoding/json"
	"github.com/streadway/amqp"
)

type Repository struct {
	channel *amqp.Channel
}

func NewRepository(ch *amqp.Channel) *Repository {
	return &Repository{channel: ch}
}

func (r *Repository) SendRabbitMessage(content *model.Metadata) error {
	jsonContent, err := json.Marshal(content)
	if err != nil {
		return err
	}

	err = r.channel.Publish("content", "file", false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         jsonContent,
		})

	if err != nil {
		return err
	}

	return nil
}
