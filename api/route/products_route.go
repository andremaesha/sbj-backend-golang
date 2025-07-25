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

// SetupPublicProductRoutes sets up routes that are publicly accessible
func SetupPublicProductRoutes(env *bootstrap.Env, session *session.Store, timeout time.Duration, db *gorm.DB, redis *redis.Client, f fiber.Router) {
	pr := repository.NewProductsRepository(db, domain.TableProducts)
	ur := repository.NewUserRepository(db, redis, domain.TableUser, "session:")
	ir := repository.NewImagesRepository(db, domain.TableImages)
	pu := usecase.NewProductsUsecase(pr, ur, ir, timeout)
	pc := controller.ProductsController{
		ProductsUsecase: pu,
		Env:             env,
		Session:         session,
	}

	// Public routes - available to all users
	f.Get("/product", pc.Product)
	f.Get("/products", pc.Products)
}

// SetupAdminProductRoutes sets up routes that require admin role
func SetupAdminProductRoutes(env *bootstrap.Env, session *session.Store, timeout time.Duration, db *gorm.DB, redis *redis.Client, f fiber.Router) {
	pr := repository.NewProductsRepository(db, domain.TableProducts)
	ur := repository.NewUserRepository(db, redis, domain.TableUser, "session:")
	ir := repository.NewImagesRepository(db, domain.TableImages)
	pu := usecase.NewProductsUsecase(pr, ur, ir, timeout)
	pc := controller.ProductsController{
		ProductsUsecase: pu,
		Env:             env,
		Session:         session,
	}

	// Admin-only routes - require admin role
	f.Post("/product", pc.CreateProduct)
	f.Post("/product/images", pc.ProductImages)
}
