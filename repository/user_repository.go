package repository

import (
	"context"
	"sbj-backend/domain"
	"sbj-backend/psql"
)

type userRepository struct {
	database psql.Database
	table    string
}

func NewUserRepository(database psql.Database, table string) domain.UserRepository {
	return &userRepository{database: database, table: table}
}

func (ur *userRepository) Create(ctx context.Context, user *domain.User) error {
	table := ur.database.Table(ur.table)

	return table.InsertOne(ctx, user)
}

func (ur *userRepository) GetByEmail(c context.Context, email string) (*domain.User, error) {
	result := new(domain.User)

	err := ur.database.Table(ur.table).FindOne(c, result, "email = ?", email)
	if err != nil {
		return nil, err
	}

	return result, nil
}
