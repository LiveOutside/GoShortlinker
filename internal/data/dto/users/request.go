package users

type RegistrationRequest struct {
	Username   string `json:"username" validate:"required,min=3,max=32,alphanum"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=7"`
	RePassword string `json:"repeat_password" validate:"required,eqfield=Password"`
}

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
