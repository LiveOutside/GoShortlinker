package auth

import (
	"context"
	"errors"
	"fmt"

	genusers "goshortlinker/internal/repos/gen/users"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrSocialNotLinked = errors.New("social identity is not linked")
	ErrSocialTaken     = errors.New("social identity is linked to another user")
)

func (s *Service) issue(id int32, username string) (LoginResult, error) {
	token, expires, err := s.tokens.Issue(id, username)
	if err != nil {
		return LoginResult{}, fmt.Errorf("issue token: %w", err)
	}
	return LoginResult{UserID: id, Token: token, ExpiresAt: expires}, nil
}

func (s *Service) LoginSocial(ctx context.Context, provider, providerID string) (LoginResult, error) {
	ctx, cancel := context.WithTimeout(ctx, s.maxTimeout)
	defer cancel()

	u, err := s.queries.GetUserByIdentity(ctx, s.database, genusers.GetUserByIdentityParams{
		Provider: provider, ProviderUserID: providerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return LoginResult{}, ErrSocialNotLinked
		}
		return LoginResult{}, fmt.Errorf("get identity: %w", err)
	}
	if !u.IsActive {
		return LoginResult{UserID: u.ID}, ErrNotActivated
	}
	return s.issue(u.ID, u.Username)
}

func (s *Service) LinkSocial(ctx context.Context, userID int32, provider, providerID string) error {
	ctx, cancel := context.WithTimeout(ctx, s.maxTimeout)
	defer cancel()

	err := s.queries.UpsertIdentity(ctx, s.database, genusers.UpsertIdentityParams{
		UserID: userID, Provider: provider, ProviderUserID: providerID,
	})
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrSocialTaken
	}
	return err
}

func (s *Service) UnlinkSocial(ctx context.Context, userID int32, provider string) error {
	ctx, cancel := context.WithTimeout(ctx, s.maxTimeout)
	defer cancel()
	_, err := s.queries.DeleteIdentity(ctx, s.database, genusers.DeleteIdentityParams{
		UserID: userID, Provider: provider,
	})
	return err
}

func (s *Service) LinkedSocials(ctx context.Context, userID int32) (map[string]bool, error) {
	ctx, cancel := context.WithTimeout(ctx, s.maxTimeout)
	defer cancel()

	list, err := s.queries.ListIdentities(ctx, s.database, userID)
	if err != nil {
		return nil, err
	}
	res := make(map[string]bool, len(list))
	for _, p := range list {
		res[p] = true
	}
	return res, nil
}
