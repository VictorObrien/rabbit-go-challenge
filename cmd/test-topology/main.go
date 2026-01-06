package main

import (
	"context"
	"log"
	"time"

	"github.com/VictorObrien/rabbitmq-go-challenge/pkg/amqp"
)

func main() {
	ctx := context.Background()

	config := amqp.DefaultConnectionConfig("amqp://rabbit:rabbit@localhost:5672/")
	topologyConfig := amqp.DefaultTopologyConfig()

	conn, err := amqp.SetupWithRetry(ctx, config, topologyConfig)
	if err != nil {
		log.Fatalf("Failed to setup: %v", err)
	}
	defer conn.Close()

	log.Println("✅ Topology declared successfully!")
	log.Println("Check RabbitMQ Management UI at http://localhost:15672")
	log.Println("  User: rabbit")
	log.Println("  Pass: rabbit")

	time.Sleep(2 * time.Second)
}


