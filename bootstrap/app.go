package bootstrap

import (
	"sbj-backend/psql"
)

type Application struct {
	Env  *Env
	Psql psql.Client
}

func App() *Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Psql = NewPsql(app.Env)

	return app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Psql)
}
