package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sbj-backend/bootstrap"
	"time"
)

func Setup(env *bootstrap.Env, session *session.Store, timeout time.Duration, db *gorm.DB, redis *redis.Client, f *fiber.App) {
	publicRouter := f.Group("/api/v1")
	NewSignupRouter(env, timeout, db, redis, publicRouter)
	NewLoginRouter(env, session, timeout, db, redis, publicRouter)
	NewLogoutRouter(env, session, timeout, db, redis, publicRouter)

	NewProductsRouter(env, session, timeout, db, redis, publicRouter)
}
