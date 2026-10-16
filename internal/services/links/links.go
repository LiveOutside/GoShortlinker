package links

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	dtolinks "goshortlinker/internal/data/dto/links"
	genlinks "goshortlinker/internal/repos/gen/links"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) GetLinkForRedirect(share_code string) (dtolinks.RedirectShareCodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.maxTimeout)
	defer cancel()

	link, err := s.queries.GetLinkByShareCode(ctx, s.database, share_code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("Link not found: %v", err)
			return dtolinks.RedirectShareCodeResponse{}, ErrLinkNotFound
		}
		log.Printf("failed to get link: %v", err)
		return dtolinks.RedirectShareCodeResponse{}, ErrInternal
	}

	if !link.IsActive {
		log.Printf("link is no longer active")
		return dtolinks.RedirectShareCodeResponse{}, ErrLinkUnavailable
	}

	if link.ValidUntil.Valid && link.ValidUntil.Time.Before(time.Now()) {
		log.Printf("link expired")
		return dtolinks.RedirectShareCodeResponse{}, ErrLinkUnavailable
	}

	if link.AllowedRedirects.Valid && link.Redirects >= link.AllowedRedirects.Int32 {
		log.Printf("link has reached redirect threshold")
		return dtolinks.RedirectShareCodeResponse{}, ErrLinkUnavailable
	}

	return dtolinks.RedirectShareCodeResponse{
		ShareCode:     link.ShareCode,
		RedirectTimer: link.RedirectTimer,
	}, nil
}

var allowedTimers = map[int32]struct{}{0: {}, 3: {}, 5: {}, 10: {}, 15: {}}

func (s *Service) SaveAndShortenLink(userID pgtype.Int4, request dtolinks.SaveRequest) (dtolinks.LinkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.maxTimeout)
	defer cancel()

	params := genlinks.SaveLinkParams{
		CreatedBy:           userID,
		ShareCode:           GenerateSharecode(),
		RedirectTimer:       5,
		RedirectTo:          request.RedirectTo,
		ValidUntil:          request.ValidUntil,
		OnlyUniqueRedirects: false,
		IsActive:            true,
	}

	if userID.Valid {
		if _, ok := allowedTimers[request.RedirectTimer]; !ok {
			return dtolinks.LinkResponse{}, ErrInvalidParams
		}
		if request.AllowedRedirects.Valid && request.AllowedRedirects.Int32 < 1 {
			return dtolinks.LinkResponse{}, ErrInvalidParams
		}

		params.RedirectTimer = request.RedirectTimer
		params.AllowedRedirects = request.AllowedRedirects
		params.OnlyUniqueRedirects = request.OnlyUniqueRedirects
		params.IsActive = request.IsActive
	}

	link, err := s.queries.SaveLink(ctx, s.database, params)
	if err != nil {
		log.Printf("Failed to shorten link: %v", err)
		return dtolinks.LinkResponse{}, ErrSaveAndShortenLink
	}

	// REWORK: change domain through .env when Docker added
	return dtolinks.LinkResponse{
		ShareCode: link.ShareCode,
	}, nil
}

func (s *Service) ConfirmRedirect(shareCode, ip string) (dtolinks.RedirectResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.maxTimeout)
	defer cancel()

	link, err := s.queries.GetLinkByShareCode(ctx, s.database, shareCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("Link not found: %v", err)
			return dtolinks.RedirectResponse{}, ErrLinkNotFound
		}
		log.Printf("failed to get link: %v", err)
		return dtolinks.RedirectResponse{}, ErrInternal
	}

	if !link.OnlyUniqueRedirects {
		return s.countRedirect(ctx, s.database, link.ID)
	}

	tx, err := s.database.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return dtolinks.RedirectResponse{}, ErrInternal
	}
	defer tx.Rollback(ctx)

	inserted, err := s.queries.RecordVisit(ctx, tx, genlinks.RecordVisitParams{
		LinkID: link.ID,
		IpHash: s.hashIP(ip),
	})
	if err != nil {
		log.Printf("failed to record visit: %v", err)
		return dtolinks.RedirectResponse{}, ErrInternal
	}

	if inserted == 0 {
		if !link.IsActive || (link.ValidUntil.Valid && link.ValidUntil.Time.Before(time.Now())) {
			return dtolinks.RedirectResponse{}, ErrLinkUnavailable
		}
		return dtolinks.RedirectResponse{RedirectTo: link.RedirectTo}, nil
	}

	response, err := s.countRedirect(ctx, tx, link.ID)
	if err != nil {
		return dtolinks.RedirectResponse{}, ErrInternal
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("failed to commit transaction: %v", err)
		return dtolinks.RedirectResponse{}, ErrInternal
	}

	return response, nil
}

func (s *Service) countRedirect(ctx context.Context, db genlinks.DBTX, linkID int32) (dtolinks.RedirectResponse, error) {
	row, err := s.queries.ConfirmRedirect(ctx, db, linkID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dtolinks.RedirectResponse{}, ErrLinkUnavailable
		}
		log.Printf("failed to confirm redirect: %v", err)
		return dtolinks.RedirectResponse{}, ErrInternal
	}
	return dtolinks.RedirectResponse{
		RedirectTo: row.RedirectTo,
	}, nil
}

func (s *Service) hashIP(ip string) string {
	mac := hmac.New(sha256.New, s.ipHashKey)
	mac.Write([]byte(ip))
	return hex.EncodeToString(mac.Sum(nil))
}
