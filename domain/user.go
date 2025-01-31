package domain

import (
	"context"
	"time"
)

const TableUser = "dbo.users"

type User struct {
	Id          int        `gorm:"primaryKey;autoIncrement" redis:"user_id"`
	Verified    bool       `gorm:"default:false"`
	FirstName   string     `gorm:"not null"`
	LastName    string     `gorm:"not null"`
	Email       string     `gorm:"unique;not null" redis:"email"`
	Password    string     `gorm:"not null"`
	AvatarId    int        `gorm:"default:NULL"`
	Role        string     `gorm:"default:'user'" redis:"role"`
	CreatedDate *time.Time `gorm:"default:NULL"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(c context.Context, email string) (*User, error)
	SetExpire(expire int)
	SetSession(ctx context.Context, idSession string, user *User) error
}
