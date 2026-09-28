package models

import "time"

type ActivationCodes struct {
	ID        uint
	UserID    *uint
	Code      string
	ExpiresAt time.Time
	CreatedAt time.Time
}
