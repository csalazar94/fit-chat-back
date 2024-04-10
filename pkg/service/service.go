package service

import (
	"context"

	"github.com/csalazar94/fit-chat-back/internal/db"
)

type IUserService interface {
	Create(context.Context, db.CreateUserParams) (User, error)
	GetAll(context.Context) ([]User, error)
}

type Service struct {
	UserService IUserService
}

func NewService(dbQueries *db.Queries) *Service {
	return &Service{
		UserService: NewUserService(dbQueries),
	}
}
