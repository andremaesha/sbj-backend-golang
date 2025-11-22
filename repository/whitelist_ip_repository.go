package repository

import (
	"context"
	"sbj-backend/domain"

	"gorm.io/gorm"
)

type whitelistIpRepository struct {
	db    *gorm.DB
	table string
}

func NewWhitelistIpRepository(db *gorm.DB, table string) domain.WhitelistIpRepository {
	return &whitelistIpRepository{db: db, table: table}
}

func (repo *whitelistIpRepository) GetDataByIp(c context.Context, ip string) (*domain.WhitelistIp, error) {
	result := new(domain.WhitelistIp)

	err := repo.db.WithContext(c).Table(repo.table).Where("ip = ?", ip).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
