package links

import (
	"context"
	"log"
	"os"
	"time"
	_ "time/tzdata"

	dtolinks "goshortlinker/internal/data/dto/links"
	genlinks "goshortlinker/internal/repos/gen/links"
)

var statsTZ = loadStatsTZ()

func loadStatsTZ() string {
	name := os.Getenv("STATS_TZ")
	if name == "" {
		return "UTC"
	}
	if _, err := time.LoadLocation(name); err != nil {
		log.Printf("invalid STATS_TZ %q, using UTC: %v", name, err)
		return "UTC"
	}
	return name
}

func (s *Service) GetStats(ctx context.Context, userID int32) (dtolinks.StatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.maxTimeout)
	defer cancel()

	row, err := s.queries.GetRedirectStats(ctx, s.database, genlinks.GetRedirectStatsParams{
		UserID: userID,
		Tz:     statsTZ,
	})
	if err != nil {
		log.Printf("failed to get redirect stats: %v", err)
		return dtolinks.StatsResponse{}, ErrInternal
	}

	return dtolinks.StatsResponse{
		Total:     row.Total,
		TotalMine: row.TotalMine,
		Today:     row.Today,
		TodayMine: row.TodayMine,
	}, nil
}
