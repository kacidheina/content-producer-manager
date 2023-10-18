package api

import (
	"content-producer-manager/internal"
	"content-producer-manager/internal/repository/consumer"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type ConsumerHandler struct {
	consumerRepo consumer.Repository
	router       *mux.Router
}

func NewConsumerHandler(repository consumer.Repository, router *mux.Router) *ConsumerHandler {
	return &ConsumerHandler{consumerRepo: repository, router: router}
}

func (c *ConsumerHandler) Consume(w http.ResponseWriter, r *http.Request) {
	reqVar := mux.Vars(r)["sender_id"]
	senderID, _ := strconv.Atoi(reqVar)
	if senderID == 0 {
		internal.RespondWithError(w, http.StatusBadRequest, "Missing sender_id")
		return
	}

	contentResult, err := c.consumerRepo.Consume(senderID)
	if err != nil {
		log.Printf("ERROR: %v", err)
		internal.RespondWithError(w, http.StatusInternalServerError, "Something went wrong while fetching the content from DB")
		return
	}

	if contentResult == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	marshalContent, err := json.Marshal(contentResult)
	if err != nil {
		log.Printf("ERROR: %v", err)
		internal.RespondWithError(w, http.StatusInternalServerError, "Something went wrong while marshalling the content")
		return
	}

	w.Write(marshalContent)
}
