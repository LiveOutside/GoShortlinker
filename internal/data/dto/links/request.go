package links

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type SaveRequest struct {
	RedirectTimer       int32              `json:"redirect_timer"`
	RedirectTo          string             `json:"redirect_to"`
	ValidUntil          pgtype.Timestamptz `json:"valid_until"`
	AllowedRedirects    pgtype.Int4        `json:"allowed_redirects,omitempty"`
	OnlyUniqueRedirects bool               `json:"only_unique_redirects"`
	IsActive            bool               `json:"is_active"`
}
