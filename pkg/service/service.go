package service

import (
	"context"

	"github.com/csalazar94/fit-chat-back/internal/db"
)

type iUserService interface {
	Create(context.Context, CreateUserParams) (User, error)
	GetAll(context.Context) ([]User, error)
}

type iAuthService interface {
	Login(context.Context, string, string) (bool, error)
}

type Service struct {
	UserService iUserService
	AuthService iAuthService
}

func NewService(dbQueries *db.Queries) *Service {
	return &Service{
		UserService: NewUserService(dbQueries),
		AuthService: NewAuthService(dbQueries),
	}
}
