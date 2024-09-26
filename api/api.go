package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/astoyanov87/subscription-service/redis"
)

type SubscriptionRequest struct {
	Email   string `json:"email"`
	MatchID string `json:"match_id"`
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	var subReq SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&subReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := redis.AddSubscriber(subReq.MatchID, subReq.Email)
	if err != nil {
		http.Error(w, "Error subscribing", http.StatusInternalServerError)
		return
	}
	fmt.Println("Subscription stored successfully in Redis")

	w.WriteHeader(http.StatusCreated)
}
