package rediswrap

import (
	"context"
	"fmt"
	"go-redis-chatroom/constants"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	// Client is redis client
	Client *redis.Client
)

func init() {
	var err error
	Client, err = Connect()
	if err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// Connect return a redis client
func Connect() (*redis.Client, error) {
	host := getEnv("HOST", "localhost")
	port := getEnv("PORT", "6379")

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	client.SAdd(ctx, constants.ChannelsKey, "general")
	return client, nil
}
