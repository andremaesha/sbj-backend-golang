package domain

import (
	"context"
	"database/sql"
)

const TableImages = "dbo.images"

type Images struct {
	Id         int           `gorm:"primaryKey;autoIncrement"`
	Url        string        `gorm:"not null"`
	CreatedBy  string        `gorm:"default:'SYSTEM'"`
	CreatedAt  *sql.NullTime `gorm:"default:now()"`
	UpdatedBy  string        `gorm:"default:null"`
	UpdatedAt  *sql.NullTime `gorm:"default:null"`
	ProductsId int           `gorm:"not null"`
}

type ImagesRepository interface {
	GetDataByProductsId(ctx context.Context, id int) []*Images
	Create(ctx context.Context, image *Images) error
}
