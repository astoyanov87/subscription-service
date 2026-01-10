package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/astoyanov87/subscription-service/config"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func AddSubscriber(matchID, email string, cfg *config.Config) error {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Host + ":" + cfg.Redis.Port,
	})

	err := rdb.SAdd(ctx, "subscribers:"+matchID, email).Err()
	if err != nil {
		log.Printf("Error adding subscriber to Redis: %v", err)
		return err
	}
	return nil
}

func GetSubscribers(matchID string, cfg *config.Config) ([]string, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Host + ":" + cfg.Redis.Port,
		DB:   0,
	})

	fmt.Println("Match ID is: subscribers:" + matchID)
	subscribers, err := rdb.SMembers(ctx, "subscribers:"+matchID).Result()
	if err != nil {
		log.Printf("Error getting subscribers from Redis: %v", err)
		return nil, err
	}

	return subscribers, nil
}
