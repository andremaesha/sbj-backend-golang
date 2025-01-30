package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"sbj-backend/api/controller"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/psql"
	"sbj-backend/redis"
	"sbj-backend/repository"
	"sbj-backend/usecase"
	"time"
)

func NewSignupRouter(env *bootstrap.Env, session *session.Store, timeout time.Duration, db psql.Database, redis redis.Database, f fiber.Router) {
	ur := repository.NewUserRepository(db, redis, domain.TableUser, env.RedisDB)
	sc := controller.SignupController{
		SignupUsecase: usecase.NewSignupUsecase(ur, timeout),
		Env:           env,
	}

	f.Post("/signup", sc.Signup)
}
