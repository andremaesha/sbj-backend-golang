package bootstrap

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(env *Env) *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addr := env.RedisHost + ":" + env.RedisPort
	username := env.RedisUser
	password := env.RedisPassword
	db := env.RedisDB

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	return client
}

func CloseRedisConnection(client *redis.Client) {
	if client == nil {
		return
	}

	err := client.Close()
	if err != nil {
		panic(err)
	}

	log.Println("Connection to redis closed.")
}
