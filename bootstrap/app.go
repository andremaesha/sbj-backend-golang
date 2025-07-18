package bootstrap

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	Env   *Env
	DB    *gorm.DB
	Redis *redis.Client
}

func App() *Application {
	app := &Application{}
	app.Env = NewEnv()
	app.DB = NewPsql(app.Env)
	app.Redis = NewRedisClient(app.Env)

	return app
}

func (app *Application) CloseDBConnection() {
	ClosePsqlConnection(app.DB)
	CloseRedisConnection(app.Redis)
}
