package users

import (
	genusers "goshortlinker/internal/repos/gen/users"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

type Service struct {
	database   genusers.DBTX
	queries    genusers.Querier
	maxTimeout time.Duration
}

func NewService(db genusers.DBTX, queries genusers.Querier, timeout time.Duration) *Service {
	return &Service{
		database:   db,
		queries:    queries,
		maxTimeout: timeout,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
