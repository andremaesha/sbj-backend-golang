package domain

import (
	"context"
	"time"
)

const TableUser = "dbo.users"

type User struct {
	Id          int        `gorm:"primaryKey;autoIncrement"`
	Verified    bool       `gorm:"default:false"`
	FirstName   string     `gorm:"not null"`
	LastName    string     `gorm:"not null"`
	Email       string     `gorm:"unique;not null"`
	Password    string     `gorm:"not null"`
	AvatarId    int        `gorm:"default:NULL"`
	Role        string     `gorm:"default:'user'"`
	CreatedDate *time.Time `gorm:"default:NULL"`
}

//type User struct {
//	Id          int
//	Verified    *sql.NullBool
//	FirstName   string
//	LastName    string
//	Email       string
//	Password    string
//	AvatarId    *sql.NullInt64
//	Role        *sql.NullString
//	CreatedDate *sql.NullString
//}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(c context.Context, email string) (*User, error)
}
