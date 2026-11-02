package auth

import (
	"context"
	"errors"
	"fmt"
	dtousers "goshortlinker/internal/data/dto/users"
	"goshortlinker/pkg/password"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
)

var dummyHash = sync.OnceValue(func() string {
	h, _ := password.Hash("dummy-password-for-constant-time")
	return h
})

func (s *Service) Login(ctx context.Context, req dtousers.AuthRequest) (LoginResult, error) {
	ctx, cancel := context.WithTimeout(ctx, s.maxTimeout)
	defer cancel()

	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.queries.GetUserByEmail(ctx, s.database, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			password.Check(req.Password, dummyHash())
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, fmt.Errorf("get user by email: %w", err)
	}

	if !password.Check(req.Password, user.PasswordHash) {
		return LoginResult{}, ErrInvalidCredentials
	}

	if !user.IsActive {
		return LoginResult{UserID: user.ID}, ErrNotActivated
	}

	token, expires, err := s.tokens.Issue(user.ID, user.Username)
	if err != nil {
		return LoginResult{}, fmt.Errorf("issue token: %w", err)
	}

	return LoginResult{UserID: user.ID, Token: token, ExpiresAt: expires}, nil
}
