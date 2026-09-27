package links

import "time"

type SaveRequest struct {
	RedirectTimer       int        `json:"redirect_timer"`
	RedirectTo          string     `json:"redirect_to"`
	ValidUntil          *time.Time `json:"valid_until"`
	AllowedRedirects    *uint      `json:"allowed_redirects,omitempty"`
	OnlyUniqueRedirects bool       `json:"only_unique_redirects"`
}
