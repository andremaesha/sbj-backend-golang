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

func NewProductsRouter(env *bootstrap.Env, session *session.Store, timeout time.Duration, db psql.Database, redis redis.Database, f fiber.Router) {
	pr := repository.NewProductsRepository(db, domain.TableProducts)
	ir := repository.NewImagesRepository(db, domain.TableImages)
	pu := usecase.NewProductsUsecase(pr, ir, timeout)
	pc := controller.ProductsController{
		ProductsUsecase: pu,
		Env:             env,
		Session:         session,
	}

	f.Get("/product", pc.Product)
}
