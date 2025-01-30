package bootstrap

import (
	"context"
	"log"
	"sbj-backend/redis"
	"time"
)

func NewRedis(env *Env) redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addr := env.RedisHost + ":" + env.RedisPort
	username := env.RedisUser
	password := env.RedisPassword
	db := env.RedisDB

	client, err := redis.NewClient(addr, username, password, db)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		panic(err)
	}

	return client
}

func CloseRedisConnection(client redis.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect()
	if err != nil {
		panic(err)
	}

	log.Println("Connection to redis closed.")
}
