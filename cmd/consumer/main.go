package main

import (
	"content-producer-manager/configs"
	"content-producer-manager/internal/api"
	"content-producer-manager/internal/repository/consumer"
	"content-producer-manager/internal/repository/consumer/postgres"
	"content-producer-manager/pkg/migrations"
	"content-producer-manager/service"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"time"
)

var (
	db                 *sql.DB
	err                error
	ch                 *amqp.Channel
	consumerRepository consumer.Repository
)

func main() {

	initRabbitMQ()
	initPostgresDB()
	initRepositories()
	initService()
	startServer()
}

func initPostgresDB() {
	time.Sleep(3 * time.Second)

	connString := fmt.Sprintf("postgres://%s:%s@db:5432/%s?sslmode=disable", configs.DbUser, configs.DbPassword, configs.DbName)

	db, err = sql.Open("postgres", connString)
	if err != nil {
		log.Print(err.Error())
		panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Print(err.Error())
		panic(err)
	}

	fmt.Println("Connected to the PostgresSQL database.")

	//Run database connection
	err = migrations.Up(configs.DbMigrationsDir, connString)

	if err != nil {
		log.Print(err.Error())
		panic(err)
	}

	fmt.Println("Database migration applied.")
}

func initRepositories() {
	consumerRepository = postgres.NewRepository(db)
}

func initService() {
	err := service.NewService(consumerRepository, ch).ConsumeMessages()
	if err != nil {
		log.Fatalf("failed to consume messages, err: %s", err.Error())
	}
}

func startServer() {
	router := mux.NewRouter()
	h := api.NewConsumerHandler(consumerRepository, router)

	subrouter := router.PathPrefix("/api").Subrouter()
	subrouter.HandleFunc("/consumer/{sender_id}", h.Consume).Methods("GET")

	err = http.ListenAndServe(fmt.Sprintf(":%d", configs.ConsumerPort), router)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println("server is listening on :8001")
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

		_, err = ch.QueueDeclare("store-files", true, false, false, false, nil)
		if err != nil {
			log.Fatalf("Failed to declare a queue: %v", err)
		}

		err = ch.QueueBind("store-files", "file", "content", false, nil)
		if err != nil {
			log.Fatalf("Failed to bind a queue: %v", err)
		}
		return
	}

	log.Println("Failed to connect to channel, retrying...", err)

}
