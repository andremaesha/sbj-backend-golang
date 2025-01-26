package bootstrap

import (
	"context"
	"log"
	"sbj-backend/psql"
	"time"
)

func NewPsql(env *Env) psql.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := dsn(
		env.DBHost,
		env.DBPort,
		env.DBUser,
		env.DBPass,
		env.DBName,
	)

	client, err := psql.NewClient(config)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		panic(err)
	}

	return client
}

func CloseMongoDBConnection(client psql.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect()
	if err != nil {
		panic(err)
	}

	log.Println("Connection to MongoDB closed.")
}
