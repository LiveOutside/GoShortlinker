package activationcodes

type ActivateRequest struct {
	UserID int32  `json:"user_id" validate:"required"`
	Code   string `json:"code" validate:"required,len=6,numeric"`
}

type ResendRequest struct {
	UserID int32 `json:"user_id" validate:"required"`
}
