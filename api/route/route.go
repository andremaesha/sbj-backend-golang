package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/internal/middlewares"
	"sbj-backend/repository"
	"sbj-backend/usecase"
	"time"
)

func Setup(env *bootstrap.Env, session *session.Store, timeout time.Duration, db *gorm.DB, redis *redis.Client, f *fiber.App) {
	// Initialize user repository for auth middleware
	userRepo := repository.NewUserRepository(db, redis, domain.TableUser, "session:")
	lookupRepository := repository.NewReffLookupRepository(db, domain.TableReffLookup)
	whitelistIpRepository := repository.NewWhitelistIpRepository(db, domain.TableWhitelistIP)

	// Initialize auth usecase
	authUsecase := usecase.NewAuthUsecase(userRepo, lookupRepository, whitelistIpRepository, timeout)

	// Public routes - no authentication required
	publicRouter := f.Group("/api/v1")
	NewSignupRouter(env, timeout, db, redis, publicRouter)
	NewLoginRouter(env, session, timeout, db, redis, publicRouter)
	NewLogoutRouter(env, session, timeout, db, redis, publicRouter)

	// Setup public product routes
	SetupPublicProductRoutes(env, session, timeout, db, redis, publicRouter)

	// Protected routes - authentication required
	protectedRouter := f.Group("/api/v1")
	protectedRouter.Use(middlewares.AuthMiddleware(env, session, authUsecase))

	// Admin routes - admin role required
	adminRouter := protectedRouter.Group("")
	adminRouter.Use(middlewares.AdminRoleMiddleware())

	// Setup admin product routes
	SetupAdminProductRoutes(env, session, timeout, db, redis, adminRouter)
}
