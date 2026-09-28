package links

import (
	"context"
	dtolinks "goshortlinker/internal/data/dto/links"
	genlinks "goshortlinker/internal/repos/gen/links"
	"log"
)

func (s *Service) SaveAndShortenLink(request dtolinks.SaveRequest) (dtolinks.LinkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.maxTimeout)
	defer cancel()

	params := genlinks.SaveLinkParams{
		// CreatedBy: userID,
		ShareCode:           GenerateSharecode(),
		RedirectTimer:       5,
		RedirectTo:          request.RedirectTo,
		ValidUntil:          request.ValidUntil,
		OnlyUniqueRedirects: false,
		IsActive:            true,
	}

	// if userID.Valid { // realisation of some kind of validation that user is auth
	// 	params.RedirectTimer = request.RedirectTimer
	// 	params.AllowedRedirects = request.AllowedRedirects
	// 	params.OnlyUniqueRedirects = request.OnlyUniqueRedirects
	// 	params.IsActive = request.IsActive
	// }

	link, err := s.queries.SaveLink(ctx, s.database, params)

	if err != nil {
		log.Printf("Failed to shorten link: %v", err)
		return dtolinks.LinkResponse{}, ErrSaveAndShortenLink
	}

	// REWORK: change domain through .env when Docker added
	return dtolinks.LinkResponse{
		ShortenedLink: "127.0.0.1/" + link.ShareCode,
	}, nil

	// return dtolinks.SaveResponse{
	// 	ID:            link.ID,
	// 	ShareCode:     link.ShareCode,
	// 	RedirectTimer: link.RedirectTimer,
	// 	RedirectTo:    link.RedirectTo,
	// 	ValidUntil:    link.ValidUntil,
	// 	IsActive:      link.IsActive,
	// 	DateCreated:   link.DateCreated,
	// }, nil
}
