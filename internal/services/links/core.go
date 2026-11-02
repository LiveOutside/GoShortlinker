package links

import (
	"context"
	"fmt"
	genlinks "goshortlinker/internal/repos/gen/links"
	"time"

	"github.com/jackc/pgx/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const shareCodeAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

type DB interface {
	genlinks.DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Service struct {
	database   DB
	queries    genlinks.Querier
	maxTimeout time.Duration
	ipHashKey  []byte
}

func NewService(db DB, queries genlinks.Querier, timeout time.Duration, ipHashKey []byte) *Service {
	return &Service{
		database:   db,
		queries:    queries,
		maxTimeout: timeout,
		ipHashKey:  ipHashKey,
	}

}

func GenerateSharecode() string {
	code, err := gonanoid.Generate(shareCodeAlphabet, 10)
	if err != nil {
		panic(fmt.Errorf("failed to generate sharecode: %w", err))
	}

	return code
}
