package main

import (
	"content-producer-manager/configs"
	"content-producer-manager/internal/api"
	"content-producer-manager/internal/repository/sender/rabbitmq"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"time"
)

var (
	ch               *amqp.Channel
	senderRepository *rabbitmq.Repository
)

func main() {
	initRabbitMQ()
	initRepositories()
	startServer()
}

func initRepositories() {
	senderRepository = rabbitmq.NewRepository(ch)
}

func startServer() {
	router := mux.NewRouter()
	h := api.NewSenderHandler(router, senderRepository)

	subrouter := router.PathPrefix("/api/").Subrouter()
	subrouter.HandleFunc("/send", h.Send).Methods("POST")

	err := http.ListenAndServe(fmt.Sprintf(":%d", configs.ApiPort), router)
	if err != nil {
		fmt.Print(err.Error())
	}

	fmt.Println("service is listening on :8002")
}

func initRabbitMQ() {

	time.Sleep(25 * time.Second)
	log.Println("Connecting to: amqp://user:password123@rabbitmq:5672/")

	// If connection is successful, return new instance
	conn, err := amqp.Dial("amqp://user:password123@rabbitmq:5672/")

	if err == nil {
		log.Println("Successfully connected to queue!")
		ch, err = conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open a channel: %v", err)
		}

		err = ch.ExchangeDeclare("content", "topic", true, false, false, false, nil)
		if err != nil {
			log.Fatalf("Failed to declare an exchange: %v", err)
		}
		fmt.Println("Exchange created.")
		return
	}

	log.Println("Failed to connect to queue, retrying...", err)

}
