package repository

import (
	"context"
	"sbj-backend/domain"
	"sbj-backend/psql"
)

type imagesRepository struct {
	database psql.Database
	table    string
}

func NewImagesRepository(database psql.Database, table string) domain.ImagesRepository {
	return &imagesRepository{database: database, table: table}
}

func (ir *imagesRepository) GetDataById(ctx context.Context, id int) *domain.Images {
	result := new(domain.Images)

	err := ir.database.Table(ir.table).FindOne(ctx, result, "id = ?", id)
	if err != nil {
		panic(err)
	}

	return result
}

func (ir *imagesRepository) Create(ctx context.Context, image *domain.Images) error {
	err := ir.database.Table(ir.table).InsertOne(ctx, image)
	if err != nil {
		return err
	}

	return nil
}
