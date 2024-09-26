package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/astoyanov87/subscription-service/email"
	"github.com/astoyanov87/subscription-service/redis"

	"github.com/streadway/amqp"
)

type MatchStatusChangedEvent struct {
	MatchID         string `json:"matchID"`
	Status          string `json:"status"`
	Name            string `json:"matchName"`
	HomePlayerScore int    `json:"homePlayerScore"`
	AwayPlayerScore int    `json:"awayPlayerScore"`
}

func ListenForMatchEvents() {
	conn, err := amqp.Dial("amqp://guest:guest@10.133.66.153:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"match_status_queue", // Queue name
		true,                 // Durable
		false,                // Delete when unused
		false,                // Exclusive
		false,                // No-wait
		nil,                  // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Bind the queue to the exchange
	err = ch.QueueBind(
		q.Name,                  // queue name
		"",                      // routing key (not used for fanout)
		"match_status_exchange", // exchange name
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer tag
		true,   // Auto-ack
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		var event MatchStatusChangedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Error decoding message: %v", err)
			continue
		}

		// Get subscribers from Redis and send emails
		subscribers, err := redis.GetSubscribers(event.MatchID)
		fmt.Println(subscribers)
		if err != nil {
			log.Printf("Error retrieving subscribers: %v", err)
			continue
		}

		for _, emailAddr := range subscribers {
			fmt.Println("Sending email to " + emailAddr)
			email.SendEmail(emailAddr, event.Name, event.HomePlayerScore, event.AwayPlayerScore, event.Status)
		}
	}
}
