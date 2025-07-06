package psql

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"log"
	"os"
	"time"
)

type Database interface {
	Table(string) Table
	Client() Client
}

type Table interface {
	FindOne(context.Context, any, string, ...any) error
	InsertOne(context.Context, any) error
	InsertMany(context.Context, []any) error
	DeleteOne(context.Context, string, ...any) (int64, error)
	Find(context.Context, any, string, ...any) error
	CountDocuments(context.Context, any, string, ...any) (int64, error)
	UpdateOne(context.Context, string, any, ...any) error
	UpdateMany(context.Context, string, any, ...any) error
}

type Client interface {
	Database() Database
	Disconnect() error
	Ping(context.Context) error
}

type gormClient struct {
	db *gorm.DB
}

type gormDatabase struct {
	db *gorm.DB
}

type gormTable struct {
	db    *gorm.DB
	table string
}

func NewClient(dsn string) (Client, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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
			SetMaxOpenConns(50),
	)
	if err != nil {
		panic(err)
	}

	return &gormClient{db: db}, nil
}

func (gc *gormClient) Disconnect() error {
	db, err := gc.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func (gc *gormClient) Ping(ctx context.Context) error {
	db, err := gc.db.DB()
	if err != nil {
		return err
	}

	return db.PingContext(ctx)
}

func (gc *gormClient) Database() Database {
	return &gormDatabase{db: gc.db}
}

func (gd *gormDatabase) Table(tableName string) Table {
	return &gormTable{db: gd.db, table: tableName}
}

func (gd *gormDatabase) Client() Client {
	return &gormClient{db: gd.db}
}

func (gt *gormTable) FindOne(ctx context.Context, result any, query string, args ...any) error {
	return gt.db.WithContext(ctx).Table(gt.table).Where(query, args...).First(result).Error
}

func (gt *gormTable) InsertOne(ctx context.Context, record any) error {
	return gt.db.WithContext(ctx).Table(gt.table).Create(record).Error
}

func (gt *gormTable) InsertMany(ctx context.Context, records []any) error {
	return gt.db.WithContext(ctx).Table(gt.table).Create(&records).Error
}

func (gt *gormTable) DeleteOne(ctx context.Context, query string, args ...any) (int64, error) {
	result := gt.db.WithContext(ctx).Table(gt.table).Where(query, args...).Delete(nil)
	return result.RowsAffected, result.Error
}

func (gt *gormTable) Find(ctx context.Context, result any, query string, args ...any) error {
	return gt.db.WithContext(ctx).Table(gt.table).Where(query, args...).Find(result).Error
}

func (gt *gormTable) CountDocuments(ctx context.Context, result any, query string, args ...any) (int64, error) {
	var count int64
	err := gt.db.WithContext(ctx).Table(gt.table).Where(query, args...).Find(&result).Count(&count).Error
	return count, err
}

func (gt *gormTable) UpdateOne(ctx context.Context, query string, updates any, args ...any) error {
	return gt.db.WithContext(ctx).Table(gt.table).Where(query, args...).Updates(updates).Error
}

func (gt *gormTable) UpdateMany(ctx context.Context, query string, updates any, args ...any) error {
	return gt.db.WithContext(ctx).Table(gt.table).Where(query, args...).Updates(updates).Error
}
