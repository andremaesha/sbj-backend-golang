package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sbj-backend/api/route"
	"sbj-backend/bootstrap"
	"sbj-backend/internal/middlewares"
	"sbj-backend/internal/validator"
	"syscall"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	db := app.DB
	redis := app.Redis
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	f := fiber.New(fiber.Config{
		AppName:      "SBJ Backend",
		ServerHeader: "backend-sbj-service",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
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
	f.Use(logger.New(logger.Config{
		TimeZone:   "Asia/Jakarta",
		TimeFormat: "2006-01-02 15:04:05.000000000",
		Format:     "Time: (${time}) | Status: (${status}) | IP: (${ip}) | Latency: (${latency}) | Method: (${method}) | Path: (${path})\n",
	}))

	// Initialize validator
	validator.Initialize()

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
