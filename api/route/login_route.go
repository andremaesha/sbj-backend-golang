package route

import (
	"github.com/gofiber/fiber/v2"
	"sbj-backend/api/controller"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/psql"
	"sbj-backend/repository"
	"sbj-backend/usecase"
	"time"
)

func NewLoginRouter(env *bootstrap.Env, timeout time.Duration, db psql.Database, f fiber.Router) {
	ur := repository.NewUserRepository(db, domain.TableUser)
	lc := controller.LoginController{
		LoginUsecase: usecase.NewLoginUsecase(ur, timeout),
		Env:          env,
	}

	f.Post("/login", lc.Login)
}
