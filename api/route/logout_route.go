package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"sbj-backend/api/controller"
	"sbj-backend/bootstrap"
	"sbj-backend/psql"
	"time"
)

func NewLogoutRouter(env *bootstrap.Env, session *session.Store, timeout time.Duration, db psql.Database, f fiber.Router) {
	lc := controller.LogoutController{
		Env:     env,
		Session: session,
	}

	f.Get("/logout", lc.Logout)
}
