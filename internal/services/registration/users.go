package users

import (
	"context"
	dtousers "goshortlinker/internal/data/dto/users"
	genusers "goshortlinker/internal/repos/gen/users"
	"log"
)

func (s *Service) PostUser(request dtousers.RegistrationRequest) (dtousers.RegistrationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.maxTimeout)
	defer cancel()

	emailTaken, err := s.queries.EmailExists(ctx, s.database, request.Email)
	if err != nil {
		log.Printf("internal error on email exists")
		return dtousers.RegistrationResponse{}, ErrInternal
	}
	if emailTaken {
		log.Printf("email already exists in database")
		return dtousers.RegistrationResponse{}, ErrEmailAlreadyExists
	}

	usernameTaken, err := s.queries.UsernameExists(ctx, s.database, request.Username)
	if err != nil {
		log.Printf("internal error on username exists")
		return dtousers.RegistrationResponse{}, ErrInternal
	}
	if usernameTaken {
		log.Printf("username already exists")
		return dtousers.RegistrationResponse{}, ErrUsernameAlreadyExists
	}

	// if request.Password != request.RePassword {
	// 	return dtousers.RegistrationResponse{}, ErrPasswordMismatch
	// }

	hash, err := HashPassword(request.Password)
	if err != nil {
		log.Printf("failed to hash password")
		return dtousers.RegistrationResponse{}, err
	}

	user, err := s.queries.CreateUser(ctx, s.database, genusers.CreateUserParams{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: hash,
	})

	return dtousers.RegistrationResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
