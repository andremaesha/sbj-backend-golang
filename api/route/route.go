package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"sbj-backend/bootstrap"
	"sbj-backend/psql"
	"time"
)

func Setup(env *bootstrap.Env, session *session.Store, timeout time.Duration, db psql.Database, f *fiber.App) {
	publicRouter := f.Group("/api")
	NewSignupRouter(env, timeout, db, publicRouter)
	NewLoginRouter(env, session, timeout, db, publicRouter)
	NewLogoutRouter(env, session, timeout, db, publicRouter)
}
