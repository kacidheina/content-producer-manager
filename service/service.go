package service

import (
	"content-producer-manager/internal/repository/consumer"
	"content-producer-manager/pkg/model"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type Service struct {
	consumerRepo consumer.Repository
	channel      *amqp.Channel
}

func NewService(consumerRepo consumer.Repository, channel *amqp.Channel) *Service {
	return &Service{consumerRepo: consumerRepo, channel: channel}
}

func (s *Service) ConsumeMessages() error {
	messages, err := s.channel.Consume("store-files", "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		var err error
		var metadata *model.Metadata
		for msg := range messages {
			err = json.Unmarshal(msg.Body, &metadata)
			if err != nil {
				s.channel.Nack(msg.DeliveryTag, false, false)
				continue
			}

			err = s.consumerRepo.StoreFile(metadata)
			if err != nil {
				log.Printf("failed to store file in DB, err: %s", err.Error())
				s.channel.Nack(msg.DeliveryTag, false, true)
				continue
			}
			log.Println("file stored in DB")
			s.channel.Ack(msg.DeliveryTag, false)
		}
	}()

	return nil
}
