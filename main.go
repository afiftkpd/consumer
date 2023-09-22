package main

import (
	"consumer/models"
	"consumer/repository/es"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/joho/godotenv"

	"github.com/segmentio/kafka-go"
)

var kafkaReader *kafka.Reader
var elasticConn *elasticsearch.TypedClient

func init() {
	fmt.Println("Init.....")
	err := godotenv.Load(".env")
	kafkaHost := os.Getenv("KAFKA_HOST")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Kafka Host: " + kafkaHost)

	config := kafka.ReaderConfig{
		Brokers:  []string{kafkaHost},
		Topic:    os.Getenv("KAFKA_TOPIC"),
		GroupID:  "my-consumer-group",
		MinBytes: 1,    // Minimum number of bytes to fetch from Kafka
		MaxBytes: 10e6, // Maximum number of bytes to fetch from Kafka
	}

	kafkaReader = kafka.NewReader(config)

	elasticConn, err = elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{os.Getenv("ES_ADDRESS")},
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Consumer Running...")
	esRepo := es.NewProductRepository(elasticConn)

	// Create a context to handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer kafkaReader.Close()

	// Create a signal channel to handle Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		for {
			fmt.Println("tt")
			select {
			case <-ctx.Done():
				fmt.Println("1")
				return
			default:
				fmt.Println("2")
				msg, err := kafkaReader.ReadMessage(ctx)
				if err != nil {
					fmt.Println("Error reading message:", err)
				}
				fmt.Printf("Received message: key=%s, value=%s\n", string(msg.Key), string(msg.Value))

				switch string(msg.Key) {
				case "delete":
					err = esRepo.Delete(ctx, string(msg.Value))
					if err != nil {
						fmt.Println("Failed to delete document elasticsearch: ", err.Error())
					}
				case "update":
					product := models.Product{}
					err := json.Unmarshal(msg.Value, &product)
					if err != nil {
						fmt.Println("Error converting string to json: ", err.Error())
					}

					err = esRepo.Update(ctx, product)
					if err != nil {
						fmt.Println("Failed to update document elasticsearch: ", err.Error())
					}
				default:
					product := models.Product{}
					err := json.Unmarshal(msg.Value, &product)
					if err != nil {
						fmt.Println("Error converting string to json: ", err.Error())
					}

					err = esRepo.Store(ctx, product)
					if err != nil {
						fmt.Println("Failed to store document elasticsearch: ", err.Error())
					}
				}
			}
		}
	}()

	<-sig
	fmt.Println("Received Ctrl+C, stopping consumer...")
}
