package domain

import (
	"context"
	"time"
)

const TableReffLookup = "dbo.reff_lookup"

type ReffLookup struct {
	Id                int        `gorm:"primaryKey;autoIncrement"`
	LookupGroup       string     `gorm:"not null"`
	LookupValue       string     `gorm:"not null"`
	LookupDescription string     `gorm:"not null"`
	CreatedBy         int        `gorm:"not null"`
	CreatedAt         *time.Time `gorm:"default:now()"`
	UpdatedBy         string     `gorm:"default:null"`
	UpdatedAt         *time.Time `gorm:"default:null"`
}

type ReffLookupRepository interface {
	GetDataByGroup(c context.Context, group string) ([]*ReffLookup, error)
}
