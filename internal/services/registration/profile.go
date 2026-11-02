package users

import (
	"context"
	"errors"
	"log"

	dtousers "goshortlinker/internal/data/dto/users"
	genusers "goshortlinker/internal/repos/gen/users"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) GetProfile(ctx context.Context, userID int32) (dtousers.ProfileResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.maxTimeout)
	defer cancel()

	p, err := s.queries.GetProfile(ctx, s.database, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dtousers.ProfileResponse{}, ErrUserNotFound
		}
		log.Printf("failed to get profile %d: %v", userID, err)
		return dtousers.ProfileResponse{}, ErrInternal
	}

	return dtousers.ProfileResponse{
		Username: p.Username,
		Email:    p.Email,
		Title:    p.Title,
		Bio:      p.Bio,
	}, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID int32, req dtousers.UpdateProfileRequest) (dtousers.ProfileResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.maxTimeout)
	defer cancel()

	p, err := s.queries.UpdateProfile(ctx, s.database, genusers.UpdateProfileParams{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
		Title:    req.Title,
		Bio:      req.Bio,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dtousers.ProfileResponse{}, ErrUserNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			switch pgErr.ConstraintName {
			case "users_username_key":
				return dtousers.ProfileResponse{}, ErrUsernameAlreadyExists
			case "users_email_key":
				return dtousers.ProfileResponse{}, ErrEmailAlreadyExists
			}
		}
		log.Printf("failed to update profile %d: %v", userID, err)
		return dtousers.ProfileResponse{}, ErrInternal
	}

	return dtousers.ProfileResponse{
		Username: p.Username,
		Email:    p.Email,
		Title:    p.Title,
		Bio:      p.Bio,
	}, nil
}
