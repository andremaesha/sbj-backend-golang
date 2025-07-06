package domain

import (
	"context"
	"database/sql"
)

const TableImages = "dbo.images"

type Images struct {
	Id        int           `gorm:"primaryKey;autoIncrement"`
	Url       string        `gorm:"not null"`
	CreatedBy string        `gorm:"default:'SYSTEM'"`
	CreatedAt *sql.NullTime `gorm:"default:now()"`
	UpdatedBy string        `gorm:"default:null"`
	UpdatedAt *sql.NullTime `gorm:"default:null"`
}

type ImagesRepository interface {
	GetDataById(ctx context.Context, id int) (*Images, error)
	Create(ctx context.Context, image *Images) error
}
