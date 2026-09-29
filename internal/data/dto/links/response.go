package links

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type SaveResponse struct {
	ID            int32              `json:"id"`
	ShareCode     string             `json:"share_code"`
	RedirectTimer int32              `json:"redirect_timer"`
	RedirectTo    string             `json:"redirect_to"`
	ValidUntil    pgtype.Timestamptz `json:"valid_until"`
	IsActive      bool               `json:"is_active"`
	DateCreated   pgtype.Timestamptz `json:"date_created"`
}

type LinkResponse struct {
	ShareCode string `json:"share_code"`
}

type RedirectResponse struct {
	RedirectTo string `json:"redirect_to"`
}
