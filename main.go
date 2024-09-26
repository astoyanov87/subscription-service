package main

import (
	"log"
	"net/http"

	"github.com/astoyanov87/subscription-service/api"
	"github.com/astoyanov87/subscription-service/rabbitmq"
)

func main() {
	// Start listening for match status events
	go rabbitmq.ListenForMatchEvents()

	// Setup REST API for subscribing
	http.HandleFunc("/subscribe", api.Subscribe)
	log.Println("Subscription service running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
