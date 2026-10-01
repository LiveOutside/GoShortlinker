package activationcodes

import (
	"context"
	"errors"
	dtocodes "goshortlinker/internal/data/dto/activation_codes"
	gencodes "goshortlinker/internal/repos/gen/activation_codes"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) IssueAndSend(ctx context.Context, userID int32, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.maxTimeout)
	defer cancel()

	code := generateCode()

	err := s.queries.UpsertActivationCode(ctx, s.database, gencodes.UpsertActivationCodeParams{
		UserID:    userID,
		CodeHash:  HashCode(code),
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(CodeTTL), Valid: true},
	})

	if err != nil {
		log.Printf("failed to upsert activation code for user %d: %v", userID, err)
		return ErrInternal
	}

	go func() {
		if err := s.mailer.SendActivationCode(email, code); err != nil {
			log.Printf("failed to send activation code to email %s: %v", email, err)
		}
	}()

	return nil
}

func (s *Service) Verify(ctx context.Context, request dtocodes.ActivateRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.maxTimeout)
	defer cancel()

	ac, err := s.queries.GetActivationCodeByUserID(ctx, s.database, request.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrCodeNotFound
		}
		log.Printf("failed to verify activation code for user %d: %v", request.UserID, err)
		return ErrInternal
	}

	if ac.ExpiresAt.Time.Before(time.Now()) {
		if err := s.Resend(ctx, request.UserID); err != nil {
			return err
		}
		return ErrCodeExpired
	}

	if ac.CodeHash != HashCode(request.Code) {
		return ErrInvalidCode
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, userID int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.maxTimeout)
	defer cancel()

	return s.queries.DeleteActivationCode(ctx, s.database, userID)
}

func (s *Service) Resend(ctx context.Context, userID int32) error {
	email, err := s.queries.GetPendingEmail(ctx, s.database, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCodeNotFound
		}
		log.Printf("failed to get email for user %d: %v", userID, err)
		return ErrInternal
	}

	return s.IssueAndSend(ctx, userID, email)
}
