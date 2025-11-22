package repository

import (
	"context"
	"sbj-backend/domain"

	"gorm.io/gorm"
)

type avatarRepository struct {
	db    *gorm.DB
	table string
}

func NewAvatarRepository(db *gorm.DB, table string) domain.AvatarRepository {
	return &avatarRepository{db: db, table: table}
}

func (repo *avatarRepository) Create(ctx context.Context, avatar *domain.Avatar) error {
	return repo.db.WithContext(ctx).Table(repo.table).Create(avatar).Error
}
