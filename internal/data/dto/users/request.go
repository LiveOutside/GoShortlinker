package users

type SaveRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=7"`
}
