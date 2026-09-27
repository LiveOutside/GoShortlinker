package links

import "time"

type SaveResponse struct {
	ID            uint       `json:"id"`
	ShareCode     string     `json:"share_code"`
	RedirectTimer int        `json:"redirect_timer"`
	RedirectTo    string     `json:"redirect_to"`
	ValidUntil    *time.Time `json:"valid_until"`
	IsActive      bool       `json:"is_active"`
	DateCreated   time.Time  `json:"date_created"`
}
