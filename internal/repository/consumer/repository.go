package consumer

import "content-producer-manager/pkg/model"

type Repository interface {
	StoreFile(metadata *model.Metadata) error
	Consume(senderID int) ([]model.Metadata, error)
}
