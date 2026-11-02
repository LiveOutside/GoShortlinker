package links

type StatsResponse struct {
	Total     int64 `json:"total"`
	TotalMine int64 `json:"total_mine"`
	Today     int64 `json:"today"`
	TodayMine int64 `json:"today_mine"`
}
