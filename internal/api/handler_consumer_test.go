package api

import (
	"bytes"
	"content-producer-manager/internal/repository/consumer/mocks"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const validRequest = `{"sender_id": "%d", "file_type": "%s", "file": "%s", "is_payable": "%t", "receiver_id": "%d"}`

func createValidRequest(senderID int, fileType string, file []byte, isPayable bool, receiverID int) string {
	return fmt.Sprintf(validRequest, senderID, fileType, file, isPayable, receiverID)
}

func TestConsumerHandler_NotFound(t *testing.T) {
	repo := &mocks.Repository{}

	req, err := http.NewRequest(http.MethodPost, "/api/consumer/1", bytes.NewBuffer([]byte(createValidRequest(1, "application/json", []byte("SGVsbG8gV29ybGQh"), true, 1))))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router := mux.NewRouter()

	consumerHandler := NewConsumerHandler(repo, router)

	router.HandleFunc("/consume/{sender_id:[0-9]+}", consumerHandler.Consume)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

}
