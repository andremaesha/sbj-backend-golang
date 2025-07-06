package domain

import "time"

const TableProducts = "dbo.products"

type Product struct {
	Id           int        `gorm:"primaryKey;autoIncrement"`
	Name         string     `gorm:"not null"`
	Price        float64    `gorm:"default:0"`
	Description  string     `gorm:"type:text"`
	Ratings      float64    `gorm:"default:null"`
	ImagesId     int        `gorm:"not null"`
	Category     string     `gorm:"default:null"`
	Stock        int        `gorm:"default:0"`
	NumOfReviews int        `gorm:"default:0"`
	ReviewsId    int        `gorm:"default:null"`
	CreatedBy    string     `gorm:"default:'SYSTEM'"`
	CreatedAt    *time.Time `gorm:"default:now()"`
	UpdatedBy    string     `gorm:"default:null"`
	UpdatedAt    *time.Time `gorm:"default:null"`
}

type ProductsRepository interface {
	GetDataById(id int) (*Product, error)
}
