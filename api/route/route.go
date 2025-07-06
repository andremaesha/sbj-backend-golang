package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"sbj-backend/bootstrap"
	"sbj-backend/psql"
	"sbj-backend/redis"
	"time"
)

func Setup(env *bootstrap.Env, session *session.Store, timeout time.Duration, db psql.Database, redis redis.Database, f *fiber.App) {
	publicRouter := f.Group("/api/v1")
	NewSignupRouter(env, session, timeout, db, redis, publicRouter)
	NewLoginRouter(env, session, timeout, db, redis, publicRouter)
	NewLogoutRouter(env, session, timeout, db, redis, publicRouter)

	NewProductsRouter(env, session, timeout, db, redis, publicRouter)
}
