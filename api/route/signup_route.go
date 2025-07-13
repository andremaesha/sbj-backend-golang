package route

import (
	"github.com/gofiber/fiber/v2"
	"sbj-backend/api/controller"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/psql"
	"sbj-backend/redis"
	"sbj-backend/repository"
	"sbj-backend/usecase"
	"time"
)

func NewSignupRouter(env *bootstrap.Env, timeout time.Duration, db psql.Database, redis redis.Database, f fiber.Router) {
	ur := repository.NewUserRepository(db, redis, domain.TableUser)
	ar := repository.NewAvatarRepository(db, domain.TableAvatar)
	sc := controller.SignupController{
		SignupUsecase: usecase.NewSignupUsecase(ur, ar, timeout),
		Env:           env,
	}

	f.Post("/signup", sc.Signup)
	f.Post("/signup/avatar", sc.UploadAvatar)
}
