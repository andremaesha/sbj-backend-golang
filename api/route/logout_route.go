package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sbj-backend/api/controller"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/repository"
	"sbj-backend/usecase"
	"time"
)

func NewLogoutRouter(env *bootstrap.Env, session *session.Store, timeout time.Duration, db *gorm.DB, redis *redis.Client, f fiber.Router) {
	ur := repository.NewUserRepository(db, redis, domain.TableUser, "session:")
	lc := controller.LogoutController{
		LogoutUsecase: usecase.NewLogoutUsecase(ur, timeout),
		Env:           env,
		Session:       session,
	}

	f.Get("/logout", lc.Logout)
}
