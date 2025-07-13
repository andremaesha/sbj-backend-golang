package repository

import (
	"context"
	"sbj-backend/domain"
	"sbj-backend/psql"
)

type avatarRepository struct {
	database psql.Database
	table    string
}

func NewAvatarRepository(database psql.Database, table string) domain.AvatarRepository {
	return &avatarRepository{database: database, table: table}
}

func (repo *avatarRepository) Create(ctx context.Context, avatar *domain.Avatar) error {
	return repo.database.Table(repo.table).InsertOne(ctx, avatar)
}
