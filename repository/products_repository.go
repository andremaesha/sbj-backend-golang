package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"sbj-backend/domain"
	errCus "sbj-backend/domain/errors"
	"sbj-backend/psql"
)

type productsRepository struct {
	database psql.Database
	table    string
}

func NewProductsRepository(database psql.Database, table string) domain.ProductsRepository {
	return &productsRepository{database: database, table: table}
}

func (pr *productsRepository) GetDataById(c context.Context, id int) (*domain.Product, error) {
	result := new(domain.Product)

	err := pr.database.Table(pr.table).FindOne(c, result, "id = ? AND is_active = ?", id, true)
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

	err := pr.database.Table(pr.table).Find(c, &results, "is_active = ?", true).Error
	if err != nil {
		panic(err())
	}

	return results
}
