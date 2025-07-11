package domain

import (
	"context"
	"time"
)

const TableAvatar = "dbo.avatar"

type Avatar struct {
	Id        int        `gorm:"primaryKey;autoIncrement"`
	UserId    int        `gorm:"not null"`
	Url       string     `gorm:"type:varchar(255);not null"`
	CreatedAt *time.Time `gorm:"default:current_timestamp()"`
	UpdatedAt *time.Time `gorm:"default:null"`
}

type AvatarRepository interface {
	Create(ctx context.Context, avatar *Avatar) error
}
