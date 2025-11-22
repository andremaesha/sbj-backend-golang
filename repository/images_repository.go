package repository

import (
	"context"
	"sbj-backend/domain"

	"gorm.io/gorm"
)

type imagesRepository struct {
	db    *gorm.DB
	table string
}

func NewImagesRepository(db *gorm.DB, table string) domain.ImagesRepository {
	return &imagesRepository{db: db, table: table}
}

func (ir *imagesRepository) GetDataByProductsId(ctx context.Context, id int) []*domain.Images {
	var results []*domain.Images

	err := ir.db.WithContext(ctx).Table(ir.table).Where("products_id = ?", id).Find(&results).Error
	if err != nil {
		panic(err)
	}

	return results
}

func (ir *imagesRepository) Create(ctx context.Context, image *domain.Images) error {
	return ir.db.WithContext(ctx).Table(ir.table).Create(image).Error
}
