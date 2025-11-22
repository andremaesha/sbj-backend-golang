package repository

import (
	"context"
	"errors"
	"sbj-backend/domain"
	errCus "sbj-backend/domain/errors"

	"gorm.io/gorm"
)

type productsRepository struct {
	db    *gorm.DB
	table string
}

func NewProductsRepository(db *gorm.DB, table string) domain.ProductsRepository {
	return &productsRepository{db: db, table: table}
}

func (pr *productsRepository) GetDataById(c context.Context, id int) (*domain.Product, error) {
	result := new(domain.Product)

	err := pr.db.WithContext(c).Table(pr.table).Where("id = ? AND is_active = ?", id, true).First(result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errCus.ErrDataNotFound
		}

		panic(err)
	}

	return result, nil
}

func (pr *productsRepository) Datas(c context.Context) []*domain.Product {
	var results []*domain.Product

	err := pr.db.WithContext(c).Table(pr.table).Where("is_active = ?", true).Find(&results).Error
	if err != nil {
		panic(err)
	}

	return results
}

func (pr *productsRepository) Create(c context.Context, product *domain.Product) error {
	return pr.db.WithContext(c).Table(pr.table).Create(product).Error
}

func (pr *productsRepository) Update(c context.Context, product *domain.Product) error {
	return pr.db.WithContext(c).Table(pr.table).Updates(product).Error
}
