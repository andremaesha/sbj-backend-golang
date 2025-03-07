package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/session"
	"log"
	"os"
	"os/signal"
	"sbj-backend/api/route"
	"sbj-backend/bootstrap"
	"sbj-backend/internal/middlewares"
	"syscall"
	"time"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	db := app.Psql.Database()
	redis := app.Redis.Database()
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	f := fiber.New(fiber.Config{
		AppName:      "SBJ Backend",
		ServerHeader: "backend-sbj-service",
	})

	store := session.New(session.Config{
		CookieHTTPOnly: true,
		CookieSecure:   true,
		CookieSameSite: "Strict",
		Expiration:     time.Minute * 10,
	})

	f.Get("/metrics", monitor.New())
	f.Use(middlewares.ResponseLogger)
	f.Use(middlewares.ErrorHandler)
	route.Setup(env, store, timeout, db, redis, f)
	f.Use(middlewares.NotFoundMiddleware)

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
