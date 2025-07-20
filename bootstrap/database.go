package bootstrap

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func NewPsql(env *Env) *gorm.DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := dsn(
		env.DBHost,
		env.DBPort,
		env.DBUser,
		env.DBPass,
		env.DBName,
	)

	db, err := gorm.Open(postgres.Open(config), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: 500 * time.Millisecond,
			},
		),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		panic(err)
	}

	err = db.Use(
		dbresolver.
			Register(dbresolver.Config{}).
			SetMaxIdleConns(10).
			SetMaxOpenConns(50).
			SetConnMaxLifetime(time.Hour).
			SetConnMaxIdleTime(30 * time.Minute),
	)
	if err != nil {
		panic(err)
	}

	// Ping the database
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = sqlDB.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	return db
}

func ClosePsqlConnection(db *gorm.DB) {
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = sqlDB.Close()
	if err != nil {
		panic(err)
	}

	log.Println("Connection to psql closed.")
}
