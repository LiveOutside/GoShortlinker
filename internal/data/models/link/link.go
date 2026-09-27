package models

import (
	models "goshortlinker/internal/data/models/users"
	"time"
)

type Link struct {
	ID                  uint
	CreatedBy           *uint
	User                *models.User
	ShareCode           string
	RedirectTimer       int
	RedirectTo          string
	ValidUntil          *time.Time
	AllowedRedirects    *uint
	Redirects           uint
	DateCreated         time.Time
	OnlyUniqueRedirects bool
	IsActive            bool
}
