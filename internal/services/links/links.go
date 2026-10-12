package links

import (
	"context"
	dtolinks "goshortlinker/internal/data/dto/links"
	genlinks "goshortlinker/internal/repos/gen/links"
	"log"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) GetLink(share_code string) (dtolinks.RedirectResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.maxTimeout)
	defer cancel()

	link, err := s.queries.GetOnlyLinkByShareCode(ctx, s.database, share_code)
	if err != nil {
		log.Printf("Failed to retrieve link from share code: %v", err)
		return dtolinks.RedirectResponse{}, err
	}

	return dtolinks.RedirectResponse{
		RedirectTo: link,
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
