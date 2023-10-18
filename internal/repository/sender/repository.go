package sender

import "content-producer-manager/pkg/model"

type Repository interface {
	SendRabbitMessage(content *model.Metadata) error
}
