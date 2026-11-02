package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	minSecretLen = 32
	tokenIssuer  = "goshortlinker"
)

type Identity struct {
	UserID   int32
	Username string
}

type claims struct {
	Username string `json:"name"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenManager(secret string, ttl time.Duration) (*TokenManager, error) {
	if len(secret) < minSecretLen {
		return nil, fmt.Errorf("JWT_SECRET must be at least %d bytes (generate: openssl rand -base64 48)", minSecretLen)
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("token ttl must be positive")
	}
	return &TokenManager{secret: []byte(secret), ttl: ttl}, nil
}

func (m *TokenManager) Issue(userID int32, username string) (string, time.Time, error) {
	now := time.Now()
	expires := now.Add(m.ttl)

	c := claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tokenIssuer,
			Subject:   strconv.FormatInt(int64(userID), 10),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expires, nil
}

func (m *TokenManager) Parse(raw string) (Identity, error) {
	c := &claims{}

	_, err := jwt.ParseWithClaims(
		raw,
		c,
		func(*jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(tokenIssuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return Identity{}, ErrInvalidToken
	}

	id, err := strconv.ParseInt(c.Subject, 10, 32)
	if err != nil || id <= 0 || c.Username == "" {
		return Identity{}, ErrInvalidToken
	}
	return Identity{UserID: int32(id), Username: c.Username}, nil
}
