package activationcodes

type ActivateRequest struct {
	UserID uint   `json:"user_id" validate:"required"`
	Code   string `json:"code" validate:"required,len=6,numeric"`
}
