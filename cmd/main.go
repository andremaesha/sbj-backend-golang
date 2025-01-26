package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"sbj-backend/api/route"
	"sbj-backend/bootstrap"
	"syscall"
	"time"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	db := app.Psql.Database()
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	f := fiber.New(fiber.Config{
		AppName:      "SBJ Backend",
		ServerHeader: "backend-sbj-service",
	})

	route.Setup(env, timeout, db, f)

	go func() {
		if err := f.Listen(":7856"); err != nil {
			panic(err)
		}
	}()

	// Setup a channel to listen for interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shut down the server gracefully.
	if err := f.ShutdownWithContext(ctx); err != nil {
		log.Panic(err)
	}

	log.Println("Server gracefully stopped")
}
