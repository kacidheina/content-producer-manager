package rabbitmq

import (
	"content-producer-manager/internal/repository/sender/mocks"
	"content-producer-manager/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendRabbitMessage_OK(t *testing.T) {
	mockRepo := mocks.Repository{}

	content := &model.Metadata{
		ReceiverID: 1,
		File:       []byte("SGVsbG8gV29ybGQh"),
		FileType:   "test",
		SenderID:   1,
	}

	mockRepo.On("SendRabbitMessage", content).Return(nil)

	// Call the SendRabbitMessage method
	err := mockRepo.SendRabbitMessage(content)

	assert.NoError(t, err)

}
