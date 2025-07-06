package repository

import (
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

func (pr *productsRepository) GetDataById(id int) (*domain.Product, error) {
	result := new(domain.Product)

	err := pr.database.Table(pr.table).FindOne(nil, result, "id = ?", id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errCus.ErrDataNotFound
		}

		panic(err)
	}

	return result, nil
}
