package links

import (
	"fmt"
	genlinks "goshortlinker/internal/repos/gen/links"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const shareCodeAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

type Service struct {
	database   genlinks.DBTX
	queries    genlinks.Querier
	maxTimeout time.Duration
}

func NewService(db genlinks.DBTX, queries genlinks.Querier, timeout time.Duration) *Service {
	return &Service{
		database:   db,
		queries:    queries,
		maxTimeout: timeout,
	}

}

func GenerateSharecode() string {
	code, err := gonanoid.Generate(shareCodeAlphabet, 10)
	if err != nil {
		panic(fmt.Errorf("failed to generate sharecode: %w", err))
	}

	return code
}
