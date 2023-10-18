package api

import (
	"content-producer-manager/internal"
	"content-producer-manager/internal/repository/sender/rabbitmq"
	"content-producer-manager/pkg/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"

	"net/http"
)

type Handler struct {
	Router     *mux.Router
	SenderRepo *rabbitmq.Repository
}

func NewSenderHandler(router *mux.Router, senderRepository *rabbitmq.Repository) *Handler {
	return &Handler{
		Router:     router,
		SenderRepo: senderRepository,
	}
}

// Send handles the send request.
func (s *Handler) Send(w http.ResponseWriter, r *http.Request) {
	var content *model.Metadata
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&content); err != nil {
		log.Printf("ERROR: %v", err)
		internal.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if content == nil || content.Validate() != nil {
		log.Printf("ERROR: %v", content.Validate())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := s.SenderRepo.SendRabbitMessage(content)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("INFO: Message sent successfully in the queue")
	w.WriteHeader(http.StatusOK)
}
