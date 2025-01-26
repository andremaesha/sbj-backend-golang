package route

import (
	"github.com/gofiber/fiber/v2"
	"sbj-backend/bootstrap"
	"sbj-backend/psql"
	"time"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db psql.Database, f *fiber.App) {
	publicRouter := f.Group("/api")
	NewSignupRouter(env, timeout, db, publicRouter)
	NewLoginRouter(env, timeout, db, publicRouter)
}
