package repository

import (
	"context"
	"sbj-backend/domain"

	"gorm.io/gorm"
)

type reffLookupRepository struct {
	db    *gorm.DB
	table string
}

func NewReffLookupRepository(db *gorm.DB, table string) domain.ReffLookupRepository {
	return &reffLookupRepository{db: db, table: table}
}

func (repo *reffLookupRepository) GetDataByGroup(c context.Context, group string) ([]*domain.ReffLookup, error) {
	results := make([]*domain.ReffLookup, 0)

	err := repo.db.WithContext(c).Table(repo.table).Where("lookup_group = ?", group).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
