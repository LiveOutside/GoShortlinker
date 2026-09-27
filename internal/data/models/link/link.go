package models

import (
	models "goshortlinker/internal/data/models/users"
	"time"
)

type Link struct {
	ID                  uint         `gorm:"primaryKey"`
	CreatedBy           *uint        `gorm:"index"`
	User                *models.User `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnDelete:SET NULL"`
	ShareCode           string       `gorm:"size:10;uniqueIndex;not null"`
	RedirectTimer       int          `gorm:"not null;default:5;check:redirect_timer IN (0,3,5,10,15)"`
	RedirectTo          string       `gorm:"not null"`
	ValidUntil          *time.Time   `gorm:"index"`
	AllowedRedirects    *uint        `gorm:"check:allowed_redirects IS NULL OR allowed_redirects >= 1"`
	Redirects           uint         `gorm:"not null;default:0"`
	DateCreated         time.Time    `gorm:"autoCreateTime"`
	OnlyUniqueRedirects bool         `gorm:"not null;default:false"`
	IsActive            bool         `gorm:"not null;default:true"`
}
