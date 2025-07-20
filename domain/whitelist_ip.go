package domain

import (
	"context"
	"time"
)

const TableWhitelistIP = "dbo.whitelist_ip"

type WhitelistIp struct {
	Id        int        `gorm:"primaryKey;autoIncrement"`
	Ip        string     `gorm:"not null;unique"`
	IsActive  bool       `gorm:"default:false"`
	CreatedBy string     `gorm:"default:'SYSTEM'"`
	CreatedAt *time.Time `gorm:"default:now()"`
	UpdatedBy string     `gorm:"default:null"`
	UpdatedAt *time.Time `gorm:"default:null"`
}

type WhitelistIpRepository interface {
	GetDataByIp(c context.Context, ip string) (*WhitelistIp, error)
}
