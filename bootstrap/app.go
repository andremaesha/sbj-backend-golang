package bootstrap

import (
	"sbj-backend/psql"
	"sbj-backend/redis"
)

type Application struct {
	Env   *Env
	Psql  psql.Client
	Redis redis.Client
}

func App() *Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Psql = NewPsql(app.Env)
	app.Redis = NewRedis(app.Env)

	return app
}

func (app *Application) CloseDBConnection() {
	ClosePsqlConnection(app.Psql)
	CloseRedisConnection(app.Redis)
}
