package activationcodes

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	gencodes "goshortlinker/internal/repos/gen/activation_codes"
	mailer "goshortlinker/pkg/mailer"
	"math/big"
	"time"
)

const CodeTTL = 15 * time.Minute

type Service struct {
	database   gencodes.DBTX
	queries    gencodes.Querier
	mailer     mailer.Mailer
	maxTimeout time.Duration
}

func NewService(db gencodes.DBTX, queries gencodes.Querier,
	mailer mailer.Mailer, timeout time.Duration) *Service {
	return &Service{
		database:   db,
		queries:    queries,
		mailer:     mailer,
		maxTimeout: timeout,
	}
}

func generateCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		panic(fmt.Errorf("failed to generate activation code: %w", err))
	}
	return fmt.Sprintf("%06d", n.Int64())
}

func HashCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}
