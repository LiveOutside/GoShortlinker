package users

import (
	genusers "goshortlinker/internal/repos/gen/users"
	activationservice "goshortlinker/internal/services/mailer"
	"time"
)

type Service struct {
	database   genusers.DBTX
	queries    genusers.Querier
	activation *activationservice.Service
	maxTimeout time.Duration
}

func NewService(db genusers.DBTX, queries genusers.Querier, activation *activationservice.Service, timeout time.Duration) *Service {
	return &Service{
		database:   db,
		queries:    queries,
		activation: activation,
		maxTimeout: timeout,
	}
}
