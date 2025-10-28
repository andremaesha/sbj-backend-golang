package domain

import (
	"context"
	"database/sql"
)

const TableImages = "dbo.images"

type Images struct {
	Id         int           `gorm:"primaryKey;autoIncrement"`
	AssetId    string        `gorm:"not null"`
	PublicId   string        `gorm:"not null"`
	Url        string        `gorm:"not null"`
	CreatedBy  int           `gorm:"not null"`
	CreatedAt  *sql.NullTime `gorm:"default:now()"`
	UpdatedBy  int           `gorm:"default:null"`
	UpdatedAt  *sql.NullTime `gorm:"default:null"`
	ProductsId int           `gorm:"not null"`
}

type ImagesRepository interface {
	GetDataByProductsId(ctx context.Context, id int) []*Images
	Create(ctx context.Context, image *Images) error
}
