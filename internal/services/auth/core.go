package auth

import (
	genusers "goshortlinker/internal/repos/gen/users"
	"time"
)

type Service struct {
	database   genusers.DBTX
	queries    genusers.Querier
	tokens     *TokenManager
	maxTimeout time.Duration
}

func NewService(db genusers.DBTX, queries genusers.Querier, tokens *TokenManager, timeout time.Duration) *Service {
	return &Service{
		database:   db,
		queries:    queries,
		tokens:     tokens,
		maxTimeout: timeout,
	}
}

type LoginResult struct {
	UserID    int32
	Token     string
	ExpiresAt time.Time
}
