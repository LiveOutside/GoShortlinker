package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"size:32;uniqueIndex;not null"`
	Email        string    `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoCreateTime"`
}
