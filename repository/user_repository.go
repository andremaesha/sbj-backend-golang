package repository

import (
	"context"
	"sbj-backend/domain"
	"sbj-backend/psql"
	"sbj-backend/redis"
	"strconv"
)

type userRepository struct {
	database   psql.Database
	redis      redis.Database
	table      string
	redisTable int
	expire     int
}

func NewUserRepository(database psql.Database, redis redis.Database, table string, redisTable int) domain.UserRepository {
	return &userRepository{database: database, redis: redis, table: table, redisTable: redisTable}
}

func (ur *userRepository) Create(ctx context.Context, user *domain.User) error {
	table := ur.database.Table(ur.table)

	return table.InsertOne(ctx, user)
}

func (ur *userRepository) SetSession(ctx context.Context, idSession string, user *domain.User) error {
	table := ur.redis.Table(strconv.Itoa(ur.redisTable))

	return table.HashSet(ctx, ur.expire, idSession, user)
}

func (ur *userRepository) GetByEmail(c context.Context, email string) (*domain.User, error) {
	result := new(domain.User)

	err := ur.database.Table(ur.table).FindOne(c, result, "email = ?", email)
	if err != nil {
		return nil, err
	}

	return result, nil
}
